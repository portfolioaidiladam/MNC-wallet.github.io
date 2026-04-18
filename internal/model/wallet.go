package model

import (
	"time"

	"github.com/google/uuid"
)

// Wallet menyimpan saldo user. Satu user tepat satu wallet (user_id UNIQUE).
// Balance dalam satuan rupiah (int64), tidak ada desimal.
type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;column:user_id"`
	Balance   int64     `gorm:"not null;default:0;column:balance"`
	CreatedAt time.Time `gorm:"not null;column:created_at"`
	UpdatedAt time.Time `gorm:"not null;column:updated_at"`
}

// TableName mengoverride nama tabel default GORM.
func (Wallet) TableName() string { return "wallets" }