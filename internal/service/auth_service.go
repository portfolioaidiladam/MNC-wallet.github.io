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
	"gorm.io/gorm"
)

// AuthService melayani register, login, dan refresh token.
type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResult, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResult, error)
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResult, error)
}

type authService struct {
	db            *gorm.DB
	users         repository.UserRepository
	wallets       repository.WalletRepository
	refreshTokens repository.RefreshTokenRepository

	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewAuthService membangun AuthService dengan dependency yang dibutuhkan.
func NewAuthService(
	db *gorm.DB,
	users repository.UserRepository,
	wallets repository.WalletRepository,
	refreshTokens repository.RefreshTokenRepository,
	jwtSecret string,
	accessTTL, refreshTTL time.Duration,
) AuthService {
	return &authService{
		db:            db,
		users:         users,
		wallets:       wallets,
		refreshTokens: refreshTokens,
		jwtSecret:     jwtSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// Register membuat user baru + wallet kosong dalam satu DB transaction.
// Phone number duplicate dimap ke ErrConflict.
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResult, error) {
	phone := strings.TrimSpace(req.PhoneNumber)
	if err := util.ValidatePhone(phone); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}
	if err := util.ValidatePIN(req.PIN); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}

	if _, err := s.users.FindByPhone(ctx, phone); err == nil {
		return nil, fmt.Errorf("%w: phone number already registered", ErrConflict)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	pinHash, err := util.HashPIN(req.PIN)
	if err != nil {
		return nil, fmt.Errorf("hash pin: %w", err)
	}

	now := time.Now()
	user := &model.User{
		ID:          model.NewID(),
		FirstName:   strings.TrimSpace(req.FirstName),
		LastName:    strings.TrimSpace(req.LastName),
		PhoneNumber: phone,
		Address:     strings.TrimSpace(req.Address),
		PINHash:     pinHash,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	wallet := &model.Wallet{
		ID:        model.NewID(),
		UserID:    user.ID,
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = repository.RunInTx(ctx, s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			if isUniqueViolation(err) {
				return fmt.Errorf("%w: phone number already registered", ErrConflict)
			}
			return fmt.Errorf("create user: %w", err)
		}
		if err := tx.WithContext(ctx).Create(wallet).Error; err != nil {
			return fmt.Errorf("create wallet: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &model.RegisterResult{
		UserID:      user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		CreatedAt:   util.FormatJakarta(user.CreatedAt),
	}, nil
}

// Login memverifikasi phone+PIN dan menerbitkan access token + refresh token.
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResult, error) {
	phone := strings.TrimSpace(req.PhoneNumber)
	if err := util.ValidatePhone(phone); err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := util.ValidatePIN(req.PIN); err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.FindByPhone(ctx, phone)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if !util.ComparePIN(user.PINHash, req.PIN) {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, user.ID)
}

// Refresh menukar refresh token yang valid dengan pasangan token baru.
// Token lama di-revoke supaya tidak bisa dipakai ulang (rotation).
func (s *authService) Refresh(ctx context.Context, refreshToken string) (*model.LoginResult, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrInvalidInput
	}
	hash := util.HashRefreshToken(refreshToken)
	rec, err := s.refreshTokens.FindByHash(ctx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	if !rec.IsActive(time.Now()) {
		return nil, ErrUnauthorized
	}

	if err := s.refreshTokens.Revoke(ctx, hash); err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, rec.UserID)
}

// issueTokens menerbitkan access token (JWT HS256) + refresh token opaque.
// Refresh token disimpan sebagai SHA-256 hex di DB.
func (s *authService) issueTokens(ctx context.Context, userID uuid.UUID) (*model.LoginResult, error) {
	access, err := util.GenerateAccessToken(s.jwtSecret, userID, s.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}
	refresh, err := util.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("issue refresh token: %w", err)
	}
	now := time.Now()
	rec := &model.RefreshToken{
		ID:        model.NewID(),
		UserID:    userID,
		TokenHash: util.HashRefreshToken(refresh),
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,
	}
	if err := s.refreshTokens.Create(ctx, rec); err != nil {
		return nil, err
	}
	return &model.LoginResult{AccessToken: access, RefreshToken: refresh}, nil
}

// isUniqueViolation mendeteksi Postgres unique constraint violation (SQLSTATE 23505).
// Pakai string match supaya tidak bergantung ke driver type.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "23505") || strings.Contains(strings.ToLower(msg), "duplicate key")
}