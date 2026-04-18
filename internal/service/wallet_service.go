package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WalletService menangani operasi saldo yang synchronous: topup dan payment.
// Transfer punya service sendiri karena butuh enqueue Asynq.
type WalletService interface {
	TopUp(ctx context.Context, userID uuid.UUID, amount int64) (*model.TopUpResult, error)
	Pay(ctx context.Context, userID uuid.UUID, amount int64, remarks string) (*model.PaymentResult, error)
}

type walletService struct {
	db           *gorm.DB
	wallets      repository.WalletRepository
	transactions repository.TransactionRepository
}

// NewWalletService membangun WalletService dengan dependency yang dibutuhkan.
func NewWalletService(
	db *gorm.DB,
	wallets repository.WalletRepository,
	transactions repository.TransactionRepository,
) WalletService {
	return &walletService{db: db, wallets: wallets, transactions: transactions}
}

// TopUp menambah saldo user dan mencatat transaction TOPUP SUCCESS dalam satu
// DB transaction dengan row-level lock di wallet.
func (s *walletService) TopUp(ctx context.Context, userID uuid.UUID, amount int64) (*model.TopUpResult, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be greater than 0", ErrInvalidInput)
	}

	var tx model.Transaction
	err := repository.RunInTx(ctx, s.db, func(gtx *gorm.DB) error {
		w, err := s.wallets.LockByUserID(ctx, gtx, userID)
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%w: wallet not found", ErrNotFound)
		}
		if err != nil {
			return err
		}
		newBalance := w.Balance + amount
		if err := s.wallets.UpdateBalance(ctx, gtx, w.ID, newBalance); err != nil {
			return err
		}
		tx = model.Transaction{
			ID:            model.NewID(),
			UserID:        userID,
			Type:          model.TxTypeTopup,
			Status:        model.TxStatusSuccess,
			Amount:        amount,
			BalanceBefore: w.Balance,
			BalanceAfter:  newBalance,
			CreatedAt:     time.Now(),
		}
		return s.transactions.Create(ctx, gtx, &tx)
	})
	if err != nil {
		return nil, err
	}

	return &model.TopUpResult{
		TopUpID:       tx.ID,
		AmountTopUp:   tx.Amount,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		CreatedAt:     util.FormatJakarta(tx.CreatedAt),
	}, nil
}

// Pay mengurangi saldo user untuk pembayaran merchant. Return ErrInsufficientFunds
// kalau saldo kurang.
func (s *walletService) Pay(ctx context.Context, userID uuid.UUID, amount int64, remarks string) (*model.PaymentResult, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be greater than 0", ErrInvalidInput)
	}

	var tx model.Transaction
	err := repository.RunInTx(ctx, s.db, func(gtx *gorm.DB) error {
		w, err := s.wallets.LockByUserID(ctx, gtx, userID)
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%w: wallet not found", ErrNotFound)
		}
		if err != nil {
			return err
		}
		if w.Balance < amount {
			return ErrInsufficientFunds
		}
		newBalance := w.Balance - amount
		if err := s.wallets.UpdateBalance(ctx, gtx, w.ID, newBalance); err != nil {
			return err
		}
		tx = model.Transaction{
			ID:            model.NewID(),
			UserID:        userID,
			Type:          model.TxTypePayment,
			Status:        model.TxStatusSuccess,
			Amount:        amount,
			BalanceBefore: w.Balance,
			BalanceAfter:  newBalance,
			Remarks:       remarks,
			CreatedAt:     time.Now(),
		}
		return s.transactions.Create(ctx, gtx, &tx)
	})
	if err != nil {
		return nil, err
	}

	return &model.PaymentResult{
		PaymentID:     tx.ID,
		Amount:        tx.Amount,
		Remarks:       tx.Remarks,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		CreatedAt:     util.FormatJakarta(tx.CreatedAt),
	}, nil
}