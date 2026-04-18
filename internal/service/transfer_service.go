package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/aidiladam/mnc-wallet/internal/worker"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// TransferService mengorkestrasi transfer antar-user: debit sender synchronous,
// lalu enqueue task Asynq untuk credit receiver.
type TransferService interface {
	Transfer(ctx context.Context, senderID, receiverID uuid.UUID, amount int64, remarks string) (*model.TransferResult, error)
}

type transferService struct {
	db           *gorm.DB
	users        repository.UserRepository
	wallets      repository.WalletRepository
	transactions repository.TransactionRepository
	asynqClient  *asynq.Client
}

// NewTransferService membangun TransferService dengan dependency yang dibutuhkan.
func NewTransferService(
	db *gorm.DB,
	users repository.UserRepository,
	wallets repository.WalletRepository,
	transactions repository.TransactionRepository,
	asynqClient *asynq.Client,
) TransferService {
	return &transferService{
		db:           db,
		users:        users,
		wallets:      wallets,
		transactions: transactions,
		asynqClient:  asynqClient,
	}
}

// Transfer melakukan debit sender synchronous dan enqueue credit async.
//
// Flow:
//  1. Validate input (amount > 0, sender != receiver, receiver exists)
//  2. Dalam DB tx: lock wallet sender, debit balance, insert TRANSFER_OUT PENDING
//  3. Enqueue task transfer:credit (max 3 retry, exponential backoff)
//  4. Return TransferResult dengan status PENDING + balance sender
//
// Kalau enqueue gagal setelah commit, kita refund sender compensating supaya
// tidak ada saldo yang "hilang". Ini jarang terjadi tapi penting.
func (s *transferService) Transfer(ctx context.Context, senderID, receiverID uuid.UUID, amount int64, remarks string) (*model.TransferResult, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be greater than 0", ErrInvalidInput)
	}
	if senderID == receiverID {
		return nil, fmt.Errorf("%w: cannot transfer to yourself", ErrInvalidInput)
	}

	if _, err := s.users.FindByID(ctx, receiverID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: target user not found", ErrNotFound)
		}
		return nil, err
	}

	var debit model.Transaction
	err := repository.RunInTx(ctx, s.db, func(gtx *gorm.DB) error {
		w, err := s.wallets.LockByUserID(ctx, gtx, senderID)
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%w: sender wallet not found", ErrNotFound)
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
		// Transaction id dipakai sebagai reference_id bersama untuk sisi IN.
		debitID := model.NewID()
		debit = model.Transaction{
			ID:                 debitID,
			UserID:             senderID,
			Type:               model.TxTypeTransferOut,
			Status:             model.TxStatusPending,
			Amount:             amount,
			BalanceBefore:      w.Balance,
			BalanceAfter:       newBalance,
			Remarks:            remarks,
			ReferenceID:        &debitID,
			CounterpartyUserID: &receiverID,
			CreatedAt:          time.Now(),
		}
		return s.transactions.Create(ctx, gtx, &debit)
	})
	if err != nil {
		return nil, err
	}

	payload := worker.TransferCreditPayload{
		ReferenceID: debit.ID,
		SenderID:    senderID,
		ReceiverID:  receiverID,
		Amount:      amount,
		Remarks:     remarks,
	}
	task, err := worker.NewTransferCreditTask(payload)
	if err != nil {
		// Refund: kembalikan saldo sender + tandai transaksi FAILED.
		s.compensate(ctx, &debit, "enqueue failed: "+err.Error())
		return nil, fmt.Errorf("build task: %w", err)
	}
	if _, err := s.asynqClient.EnqueueContext(ctx, task,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
	); err != nil {
		s.compensate(ctx, &debit, "enqueue failed: "+err.Error())
		return nil, fmt.Errorf("enqueue task: %w", err)
	}

	return &model.TransferResult{
		TransferID:    debit.ID,
		Status:        string(model.TxStatusPending),
		Amount:        debit.Amount,
		Remarks:       debit.Remarks,
		BalanceBefore: debit.BalanceBefore,
		BalanceAfter:  debit.BalanceAfter,
		CreatedAt:     util.FormatJakarta(debit.CreatedAt),
	}, nil
}

// compensate me-refund sender dan menandai transaksi debit FAILED.
// Best-effort: kalau gagal di sini, error dilog tapi tidak di-propagate
// ke caller — user akan lihat transfer gagal via /transactions (status FAILED).
func (s *transferService) compensate(ctx context.Context, debit *model.Transaction, reason string) {
	_ = repository.RunInTx(ctx, s.db, func(gtx *gorm.DB) error {
		w, err := s.wallets.LockByUserID(ctx, gtx, debit.UserID)
		if err != nil {
			return err
		}
		newBalance := w.Balance + debit.Amount
		if err := s.wallets.UpdateBalance(ctx, gtx, w.ID, newBalance); err != nil {
			return err
		}
		// Update status debit jadi FAILED via raw update (repo tidak expose).
		if err := gtx.WithContext(ctx).Model(&model.Transaction{}).
			Where("id = ?", debit.ID).
			Updates(map[string]any{
				"status":  model.TxStatusFailed,
				"remarks": truncateRemarks(debit.Remarks, reason),
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

// truncateRemarks menggabungkan remark lama + reason dengan total max 255 char.
func truncateRemarks(original, extra string) string {
	out := original
	if out != "" {
		out += " | "
	}
	out += extra
	if len(out) > 255 {
		out = out[:255]
	}
	return out
}