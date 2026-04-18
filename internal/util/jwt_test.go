package util

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseAccessToken_RoundTrip(t *testing.T) {
	secret := "unit-test-secret"
	uid := uuid.New()

	tok, err := GenerateAccessToken(secret, uid, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tok)

	claims, err := ParseAccessToken(secret, tok)
	require.NoError(t, err)
	assert.Equal(t, uid, claims.UserID)
}

func TestParseAccessToken_WrongSecretFails(t *testing.T) {
	tok, err := GenerateAccessToken("right", uuid.New(), time.Minute)
	require.NoError(t, err)

	_, err = ParseAccessToken("wrong", tok)
	assert.Error(t, err)
}

func TestParseAccessToken_ExpiredFails(t *testing.T) {
	tok, err := GenerateAccessToken("s", uuid.New(), -time.Second)
	require.NoError(t, err)

	_, err = ParseAccessToken("s", tok)
	assert.Error(t, err)
}

func TestHashRefreshToken_Deterministic(t *testing.T) {
	h1 := HashRefreshToken("abc")
	h2 := HashRefreshToken("abc")
	assert.Equal(t, h1, h2)
	assert.Len(t, h1, 64)
}

func TestValidatePhone(t *testing.T) {
	valid := []string{"0812345678", "081234567890", "0812345678901"}
	for _, v := range valid {
		assert.NoError(t, ValidatePhone(v), "expected valid: %s", v)
	}
	invalid := []string{"812345678", "091234567890", "0812abcd5678", "012345678", "08123456789012345"}
	for _, v := range invalid {
		assert.Error(t, ValidatePhone(v), "expected invalid: %s", v)
	}
}

func TestValidatePIN(t *testing.T) {
	assert.NoError(t, ValidatePIN("123456"))
	assert.Error(t, ValidatePIN("12345"))
	assert.Error(t, ValidatePIN("1234567"))
	assert.Error(t, ValidatePIN("12345a"))
}

func TestHashAndComparePIN(t *testing.T) {
	h, err := HashPIN("123456")
	require.NoError(t, err)
	assert.True(t, ComparePIN(h, "123456"))
	assert.False(t, ComparePIN(h, "654321"))
}