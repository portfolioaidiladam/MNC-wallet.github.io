package model

import (
	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

// RegisterRequest adalah payload POST /register.
type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required,max=64"`
	LastName    string `json:"last_name" binding:"required,max=64"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required,max=255"`
	PIN         string `json:"pin" binding:"required"`
}

// RegisterResult adalah payload sukses POST /register.
type RegisterResult struct {
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	CreatedAt   string    `json:"created_at"`
}

// LoginRequest adalah payload POST /login.
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PIN         string `json:"pin" binding:"required"`
}

// LoginResult adalah payload sukses POST /login.
type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshRequest adalah payload POST /refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ---------------------------------------------------------------------------
// Wallet / Transactions
// ---------------------------------------------------------------------------

// TopUpRequest adalah payload POST /topup.
type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

// TopUpResult adalah payload sukses POST /topup.
type TopUpResult struct {
	TopUpID       uuid.UUID `json:"top_up_id"`
	AmountTopUp   int64     `json:"amount_top_up"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	CreatedAt     string    `json:"created_at"`
}

// PaymentRequest adalah payload POST /pay.
type PaymentRequest struct {
	Amount  int64  `json:"amount" binding:"required,gt=0"`
	Remarks string `json:"remarks" binding:"max=255"`
}

// PaymentResult adalah payload sukses POST /pay.
type PaymentResult struct {
	PaymentID     uuid.UUID `json:"payment_id"`
	Amount        int64     `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	CreatedAt     string    `json:"created_at"`
}

// TransferRequest adalah payload POST /transfer.
type TransferRequest struct {
	TargetUser uuid.UUID `json:"target_user" binding:"required"`
	Amount     int64     `json:"amount" binding:"required,gt=0"`
	Remarks    string    `json:"remarks" binding:"max=255"`
}

// TransferResult adalah payload sukses POST /transfer.
//
// Status "PENDING" karena credit ke receiver dilakukan asynchronous oleh worker.
// Sender sudah di-debit saat response dikirim, jadi balance_before / balance_after
// sudah terisi dari sisi sender.
type TransferResult struct {
	TransferID    uuid.UUID `json:"transfer_id"`
	Status        string    `json:"status"`
	Amount        int64     `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	CreatedAt     string    `json:"created_at"`
}

// TransactionListItem merepresentasikan satu baris list GET /transactions.
// Field ID per-tipe di-populate hanya salah satu sesuai Type, sisanya omit
// supaya response rapi.
type TransactionListItem struct {
	Status          string     `json:"status"`
	UserID          uuid.UUID  `json:"user_id"`
	TransactionType string     `json:"transaction_type"`
	TopUpID         *uuid.UUID `json:"top_up_id,omitempty"`
	PaymentID       *uuid.UUID `json:"payment_id,omitempty"`
	TransferID      *uuid.UUID `json:"transfer_id,omitempty"`
	Amount          int64      `json:"amount"`
	Remarks         string     `json:"remarks"`
	BalanceBefore   int64      `json:"balance_before"`
	BalanceAfter    int64      `json:"balance_after"`
	CreatedAt       string     `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Profile
// ---------------------------------------------------------------------------

// UpdateProfileRequest adalah payload PUT /profile.
//
// phone_number & pin tidak boleh diubah di endpoint ini.
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required,max=64"`
	LastName  string `json:"last_name" binding:"required,max=64"`
	Address   string `json:"address" binding:"required,max=255"`
}

// ProfileResult adalah payload GET /profile & PUT /profile.
type ProfileResult struct {
	UserID      uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Balance     int64     `json:"balance"`
	UpdatedAt   string    `json:"updated_at"`
}