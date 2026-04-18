package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository adalah kontrak akses data tabel users.
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat UserRepository yang dibacking oleh *gorm.DB.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create insert user baru. Harus gagal (return error) kalau phone_number
// sudah ada (UNIQUE violation) — caller memetakan ke 409.
func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// FindByPhone mencari user berdasarkan phone_number. Return ErrNotFound
// kalau tidak ada.
func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("phone_number = ?", phone).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by phone: %w", err)
	}
	return &u, nil
}

// FindByID mencari user berdasarkan primary key. Return ErrNotFound
// kalau tidak ada.
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

// Update menyimpan perubahan user (first_name, last_name, address).
// phone_number & pin_hash TIDAK boleh diubah lewat endpoint /profile;
// pemanggil harus memastikan field unique key tidak berubah.
func (r *userRepository) Update(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Save(u).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}