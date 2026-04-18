package service

import (
	"context"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// UserRepository mock
// ---------------------------------------------------------------------------

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, u *model.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	v, _ := args.Get(0).(*model.User)
	return v, args.Error(1)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	v, _ := args.Get(0).(*model.User)
	return v, args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, u *model.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// ---------------------------------------------------------------------------
// WalletRepository mock
// ---------------------------------------------------------------------------

type mockWalletRepo struct {
	mock.Mock
}

func (m *mockWalletRepo) Create(ctx context.Context, w *model.Wallet) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockWalletRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	v, _ := args.Get(0).(*model.Wallet)
	return v, args.Error(1)
}

func (m *mockWalletRepo) LockByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (*model.Wallet, error) {
	args := m.Called(ctx, tx, userID)
	v, _ := args.Get(0).(*model.Wallet)
	return v, args.Error(1)
}

func (m *mockWalletRepo) UpdateBalance(ctx context.Context, tx *gorm.DB, walletID uuid.UUID, newBalance int64) error {
	args := m.Called(ctx, tx, walletID, newBalance)
	return args.Error(0)
}

// ---------------------------------------------------------------------------
// TransactionRepository mock
// ---------------------------------------------------------------------------

type mockTxRepo struct {
	mock.Mock
}

func (m *mockTxRepo) Create(ctx context.Context, tx *gorm.DB, t *model.Transaction) error {
	args := m.Called(ctx, tx, t)
	return args.Error(0)
}

func (m *mockTxRepo) ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]model.Transaction, error) {
	args := m.Called(ctx, userID, limit)
	v, _ := args.Get(0).([]model.Transaction)
	return v, args.Error(1)
}

// ---------------------------------------------------------------------------
// RefreshTokenRepository mock
// ---------------------------------------------------------------------------

type mockRefreshRepo struct {
	mock.Mock
}

func (m *mockRefreshRepo) Create(ctx context.Context, t *model.RefreshToken) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *mockRefreshRepo) FindByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	v, _ := args.Get(0).(*model.RefreshToken)
	return v, args.Error(1)
}

func (m *mockRefreshRepo) Revoke(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}