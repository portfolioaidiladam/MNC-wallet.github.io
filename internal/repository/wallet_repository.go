package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WalletRepository adalah kontrak akses data tabel wallets.
//
// Pattern row-level locking untuk operasi balance:
//
//	repository.RunInTx(ctx, db, func(tx *gorm.DB) error {
//	    w, err := walletRepo.LockByUserID(ctx, tx, userID)
//	    if err != nil { return err }
//	    newBalance := w.Balance + delta
//	    if err := walletRepo.UpdateBalance(ctx, tx, w.ID, newBalance); err != nil { return err }
//	    return txRepo.Create(ctx, tx, &model.Transaction{
//	        BalanceBefore: w.Balance, BalanceAfter: newBalance, ...
//	    })
//	})
//
// LockByUserID membuka `SELECT ... FOR UPDATE` di dalam tx, jadi concurrent
// topup/pay/transfer pada wallet yang sama akan serialized.
type WalletRepository interface {
	Create(ctx context.Context, w *model.Wallet) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Wallet, error)
	// LockByUserID melakukan SELECT ... FOR UPDATE pada wallet user
	// di dalam tx, supaya caller bisa read current balance sebelum menulis
	// balance baru. Harus dipanggil di dalam DB transaction.
	LockByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (*model.Wallet, error)
	// UpdateBalance menulis balance baru (tanpa lock ulang). Pattern-nya:
	//   w, _ := repo.LockByUserID(ctx, tx, userID)
	//   repo.UpdateBalance(ctx, tx, w.ID, w.Balance + delta)
	UpdateBalance(ctx context.Context, tx *gorm.DB, walletID uuid.UUID, newBalance int64) error
}

type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository membuat WalletRepository yang dibacking oleh *gorm.DB.
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

// Create insert wallet baru (biasanya dipanggil bersamaan dengan create user).
func (r *walletRepository) Create(ctx context.Context, w *model.Wallet) error {
	if err := r.db.WithContext(ctx).Create(w).Error; err != nil {
		return fmt.Errorf("create wallet: %w", err)
	}
	return nil
}

// FindByUserID membaca wallet tanpa lock — untuk read-only use case
// (mis. tampilkan saldo di profile). Untuk operasi write gunakan UpdateBalance.
func (r *walletRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Wallet, error) {
	var w model.Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find wallet by user: %w", err)
	}
	return &w, nil
}

// LockByUserID melakukan SELECT ... FOR UPDATE pada wallet milik userID.
// Harus dipanggil di dalam transaction. Service lalu memutuskan balance baru
// dan memanggil UpdateBalance.
func (r *walletRepository) LockByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (*model.Wallet, error) {
	if tx == nil {
		return nil, fmt.Errorf("lock wallet: tx is required")
	}
	var w model.Wallet
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("lock wallet: %w", err)
	}
	return &w, nil
}

// UpdateBalance menulis balance dan updated_at ke wallet walletID.
// Biasanya dipanggil setelah LockByUserID di dalam transaction yang sama.
func (r *walletRepository) UpdateBalance(ctx context.Context, tx *gorm.DB, walletID uuid.UUID, newBalance int64) error {
	if tx == nil {
		return fmt.Errorf("update balance: tx is required")
	}
	if err := tx.WithContext(ctx).Model(&model.Wallet{}).
		Where("id = ?", walletID).
		Updates(map[string]any{
			"balance":    newBalance,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}
	return nil
}