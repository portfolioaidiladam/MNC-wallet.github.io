package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken menyimpan hash token refresh supaya bisa di-revoke per-device.
// TokenHash = SHA-256 hex dari token plaintext (token plaintext HANYA dikirim
// ke client sekali saat login, lalu tidak disimpan).
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;column:id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index;column:user_id"`
	TokenHash string     `gorm:"size:64;not null;uniqueIndex;column:token_hash"`
	ExpiresAt time.Time  `gorm:"not null;column:expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"not null;column:created_at"`
}

// TableName mengoverride nama tabel default GORM.
func (RefreshToken) TableName() string { return "refresh_tokens" }

// IsActive mengembalikan true kalau token belum di-revoke dan belum expired.
func (t *RefreshToken) IsActive(now time.Time) bool {
	if t.RevokedAt != nil {
		return false
	}
	return now.Before(t.ExpiresAt)
}