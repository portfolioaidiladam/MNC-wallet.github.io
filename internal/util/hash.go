// Package util menyediakan helper umum: hashing, JWT, response, dll.
package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPIN membuat hash bcrypt dari PIN plaintext (cost 10).
func HashPIN(pin string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// ComparePIN membandingkan hash bcrypt dengan PIN plaintext.
func ComparePIN(hash, pin string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin)) == nil
}

// GenerateRefreshToken menghasilkan opaque refresh token 32-byte random
// (di-hex encode jadi 64 karakter). Token ini dikirim ke client sekali,
// sisi server menyimpan SHA-256 hex-nya di tabel refresh_tokens.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashRefreshToken menghasilkan SHA-256 hex dari refresh token plaintext.
// Panjang output 64 karakter — sesuai kolom refresh_tokens.token_hash.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}