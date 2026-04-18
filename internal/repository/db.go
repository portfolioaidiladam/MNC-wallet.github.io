// Package repository berisi layer akses database (GORM).
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrNotFound adalah error universal repository untuk record yang tidak ditemukan.
// Service layer membandingkan dengan errors.Is untuk memetakan ke 404.
var ErrNotFound = errors.New("record not found")

// RunInTx menjalankan fn dalam satu DB transaction. Commit kalau fn return nil,
// rollback otherwise. Context diteruskan supaya cancellation dari upstream
// diperhatikan.
func RunInTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(fn)
}