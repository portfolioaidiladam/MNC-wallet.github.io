package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"gorm.io/gorm"
)

// RefreshTokenRepository adalah kontrak akses data tabel refresh_tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, t *model.RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository membuat RefreshTokenRepository yang dibacking oleh *gorm.DB.
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, t *model.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

// FindByHash mencari refresh token berdasarkan SHA-256 hash dari token plaintext.
// Return ErrNotFound kalau tidak ada. Caller harus mengecek IsActive().
func (r *refreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var t model.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return &t, nil
}

// Revoke men-set revoked_at = now pada token dengan hash yang diberikan.
// Idempotent — revoke pada token yang sudah revoked tidak error.
func (r *refreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", now)
	if res.Error != nil {
		return fmt.Errorf("revoke refresh token: %w", res.Error)
	}
	return nil
}