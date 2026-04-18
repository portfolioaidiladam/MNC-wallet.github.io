package model

import (
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan baris di tabel users.
type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	FirstName   string    `gorm:"size:64;not null;column:first_name"`
	LastName    string    `gorm:"size:64;not null;column:last_name"`
	PhoneNumber string    `gorm:"size:16;not null;uniqueIndex;column:phone_number"`
	Address     string    `gorm:"size:255;not null;column:address"`
	PINHash     string    `gorm:"size:72;not null;column:pin_hash"`
	CreatedAt   time.Time `gorm:"not null;column:created_at"`
	UpdatedAt   time.Time `gorm:"not null;column:updated_at"`
}

// TableName mengoverride nama tabel default GORM.
func (User) TableName() string { return "users" }