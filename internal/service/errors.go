package service

import "errors"

// Error sentinel yang dipetakan handler ke HTTP status tertentu.
//
//	ErrInvalidInput       -> 400
//	ErrInvalidCredentials -> 401
//	ErrUnauthorized       -> 401
//	ErrNotFound           -> 404
//	ErrConflict           -> 409
//	ErrInsufficientFunds  -> 400 (sesuai spec: "balance not enough")
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrInsufficientFunds  = errors.New("balance not enough")
)