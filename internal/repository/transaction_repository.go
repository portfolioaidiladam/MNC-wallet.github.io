package repository

import (
	"context"
	"fmt"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionRepository adalah kontrak akses data tabel transactions.
type TransactionRepository interface {
	// Create insert transaction record. Kalau tx != nil, insert dilakukan
	// dalam transaction tersebut — dipakai saat create balance+transaction
	// harus atomic. Kalau tx == nil, pakai koneksi default repo.
	Create(ctx context.Context, tx *gorm.DB, t *model.Transaction) error

	// ListByUserID mengembalikan transaksi milik userID, di-sort DESC by created_at.
	// Dibatasi `limit` record; nilai 0 berarti tanpa batas.
	ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository membuat TransactionRepository yang dibacking oleh *gorm.DB.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, t *model.Transaction) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}
	if err := conn.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	return nil
}

func (r *transactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]model.Transaction, error) {
	q := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	var out []model.Transaction
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	return out, nil
}