// Package model berisi domain struct (GORM entities) dan DTO MNC Wallet API.
package model

import "github.com/google/uuid"

// NewID menghasilkan UUID v4 baru. Dipakai di service layer supaya repository
// tidak bergantung ke generator ID.
func NewID() uuid.UUID {
	return uuid.New()
}