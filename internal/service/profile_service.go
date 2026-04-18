package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/google/uuid"
)

// ProfileService menangani endpoint profile: GET dan PUT.
type ProfileService interface {
	Get(ctx context.Context, userID uuid.UUID) (*model.ProfileResult, error)
	Update(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.ProfileResult, error)
}

type profileService struct {
	users   repository.UserRepository
	wallets repository.WalletRepository
}

// NewProfileService membangun ProfileService.
func NewProfileService(users repository.UserRepository, wallets repository.WalletRepository) ProfileService {
	return &profileService{users: users, wallets: wallets}
}

// Get mengembalikan profil user + saldo wallet.
func (s *profileService) Get(ctx context.Context, userID uuid.UUID) (*model.ProfileResult, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	w, err := s.wallets.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return buildProfile(u, w), nil
}

// Update mengubah first_name, last_name, address. phone_number & pin tidak
// bisa diubah di endpoint ini.
func (s *profileService) Update(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.ProfileResult, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.FirstName = strings.TrimSpace(req.FirstName)
	u.LastName = strings.TrimSpace(req.LastName)
	u.Address = strings.TrimSpace(req.Address)
	u.UpdatedAt = time.Now()
	if u.FirstName == "" || u.LastName == "" || u.Address == "" {
		return nil, fmt.Errorf("%w: first_name, last_name, address cannot be empty", ErrInvalidInput)
	}
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	w, err := s.wallets.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return buildProfile(u, w), nil
}

func buildProfile(u *model.User, w *model.Wallet) *model.ProfileResult {
	return &model.ProfileResult{
		UserID:      u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
		Balance:     w.Balance,
		UpdatedAt:   util.FormatJakarta(u.UpdatedAt),
	}
}