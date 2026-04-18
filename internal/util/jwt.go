package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims adalah payload JWT access token MNC Wallet.
type Claims struct {
	UserID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

// GenerateAccessToken menerbitkan JWT HS256 dengan UserID sebagai claim `uid`.
func GenerateAccessToken(secret string, userID uuid.UUID, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAccessToken memvalidasi signature + expiry dan return claims.
func ParseAccessToken(secret, tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}