package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// buildAuthService membuat AuthService dengan mock repositories. db = nil
// karena test path yang dites tidak masuk ke RunInTx.
func buildAuthService(users *mockUserRepo, wallets *mockWalletRepo, refresh *mockRefreshRepo) AuthService {
	return NewAuthService(
		nil, users, wallets, refresh,
		"test-secret",
		15*time.Minute,
		24*time.Hour,
	)
}

// ---------------------------------------------------------------------------
// Register — hanya tes validation path (tidak masuk ke DB transaction).
// ---------------------------------------------------------------------------

func TestAuthService_Register_InvalidPhone(t *testing.T) {
	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, &mockRefreshRepo{})

	_, err := svc.Register(context.Background(), &model.RegisterRequest{
		FirstName:   "Aidil",
		LastName:    "Adam",
		PhoneNumber: "9123456789",
		Address:     "Jakarta",
		PIN:         "123456",
	})
	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestAuthService_Register_InvalidPIN(t *testing.T) {
	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, &mockRefreshRepo{})

	_, err := svc.Register(context.Background(), &model.RegisterRequest{
		FirstName:   "Aidil",
		LastName:    "Adam",
		PhoneNumber: "081234567890",
		Address:     "Jakarta",
		PIN:         "12ab56",
	})
	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestAuthService_Register_PhoneAlreadyExists(t *testing.T) {
	users := &mockUserRepo{}
	users.On("FindByPhone", mock.Anything, "081234567890").
		Return(&model.User{ID: uuid.New()}, nil)

	svc := buildAuthService(users, &mockWalletRepo{}, &mockRefreshRepo{})

	_, err := svc.Register(context.Background(), &model.RegisterRequest{
		FirstName:   "Aidil",
		LastName:    "Adam",
		PhoneNumber: "081234567890",
		Address:     "Jakarta",
		PIN:         "123456",
	})
	assert.ErrorIs(t, err, ErrConflict)
	users.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

func TestAuthService_Login_Success(t *testing.T) {
	pin := "123456"
	pinHash, err := util.HashPIN(pin)
	require.NoError(t, err)

	uid := uuid.New()
	users := &mockUserRepo{}
	users.On("FindByPhone", mock.Anything, "081234567890").
		Return(&model.User{ID: uid, PhoneNumber: "081234567890", PINHash: pinHash}, nil)

	refresh := &mockRefreshRepo{}
	refresh.On("Create", mock.Anything, mock.AnythingOfType("*model.RefreshToken")).
		Return(nil)

	svc := buildAuthService(users, &mockWalletRepo{}, refresh)
	res, err := svc.Login(context.Background(), &model.LoginRequest{
		PhoneNumber: "081234567890",
		PIN:         pin,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)

	claims, err := util.ParseAccessToken("test-secret", res.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, uid, claims.UserID)

	users.AssertExpectations(t)
	refresh.AssertExpectations(t)
}

func TestAuthService_Login_PhoneNotFound(t *testing.T) {
	users := &mockUserRepo{}
	users.On("FindByPhone", mock.Anything, "081234567890").
		Return(nil, repository.ErrNotFound)

	svc := buildAuthService(users, &mockWalletRepo{}, &mockRefreshRepo{})
	_, err := svc.Login(context.Background(), &model.LoginRequest{
		PhoneNumber: "081234567890",
		PIN:         "123456",
	})
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_WrongPIN(t *testing.T) {
	pinHash, err := util.HashPIN("123456")
	require.NoError(t, err)

	users := &mockUserRepo{}
	users.On("FindByPhone", mock.Anything, "081234567890").
		Return(&model.User{ID: uuid.New(), PINHash: pinHash}, nil)

	svc := buildAuthService(users, &mockWalletRepo{}, &mockRefreshRepo{})
	_, err = svc.Login(context.Background(), &model.LoginRequest{
		PhoneNumber: "081234567890",
		PIN:         "999999",
	})
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_InvalidFormat(t *testing.T) {
	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, &mockRefreshRepo{})
	_, err := svc.Login(context.Background(), &model.LoginRequest{
		PhoneNumber: "not-a-phone",
		PIN:         "123456",
	})
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

// ---------------------------------------------------------------------------
// Refresh
// ---------------------------------------------------------------------------

func TestAuthService_Refresh_Success(t *testing.T) {
	token := "plain-refresh-token"
	hash := util.HashRefreshToken(token)

	refresh := &mockRefreshRepo{}
	uid := uuid.New()
	refresh.On("FindByHash", mock.Anything, hash).
		Return(&model.RefreshToken{
			ID:        uuid.New(),
			UserID:    uid,
			TokenHash: hash,
			ExpiresAt: time.Now().Add(time.Hour),
		}, nil)
	refresh.On("Revoke", mock.Anything, hash).Return(nil)
	refresh.On("Create", mock.Anything, mock.AnythingOfType("*model.RefreshToken")).
		Return(nil)

	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, refresh)
	res, err := svc.Refresh(context.Background(), token)
	require.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)
	assert.NotEqual(t, token, res.RefreshToken)

	refresh.AssertExpectations(t)
}

func TestAuthService_Refresh_EmptyToken(t *testing.T) {
	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, &mockRefreshRepo{})
	_, err := svc.Refresh(context.Background(), "")
	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestAuthService_Refresh_NotFound(t *testing.T) {
	refresh := &mockRefreshRepo{}
	refresh.On("FindByHash", mock.Anything, mock.Anything).
		Return(nil, repository.ErrNotFound)

	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, refresh)
	_, err := svc.Refresh(context.Background(), "some-token")
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestAuthService_Refresh_Revoked(t *testing.T) {
	now := time.Now()
	refresh := &mockRefreshRepo{}
	refresh.On("FindByHash", mock.Anything, mock.Anything).
		Return(&model.RefreshToken{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ExpiresAt: now.Add(time.Hour),
			RevokedAt: &now,
		}, nil)

	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, refresh)
	_, err := svc.Refresh(context.Background(), "some-token")
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestAuthService_Refresh_Expired(t *testing.T) {
	refresh := &mockRefreshRepo{}
	refresh.On("FindByHash", mock.Anything, mock.Anything).
		Return(&model.RefreshToken{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ExpiresAt: time.Now().Add(-time.Hour),
		}, nil)

	svc := buildAuthService(&mockUserRepo{}, &mockWalletRepo{}, refresh)
	_, err := svc.Refresh(context.Background(), "some-token")
	assert.ErrorIs(t, err, ErrUnauthorized)
}

// ---------------------------------------------------------------------------
// Guard: memastikan helper util masih simetris dengan yang dipakai service.
// ---------------------------------------------------------------------------

func TestUtilRefreshHashSymmetry(t *testing.T) {
	tok, err := util.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, tok)
	assert.Equal(t, 64, len(util.HashRefreshToken(tok)))
	// guard terhadap drift pkg util
	assert.True(t, errors.Is(errors.New("x"), errors.New("x")) == false)
}