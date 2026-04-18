package handler

import (
	"github.com/aidiladam/mnc-wallet/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Handlers adalah bundle semua handler yang dipakai Gin engine.
type Handlers struct {
	Auth        *AuthHandler
	Wallet      *WalletHandler
	Transfer    *TransferHandler
	Transaction *TransactionHandler
	Profile     *ProfileHandler
}

// RegisterRoutes mendaftarkan seluruh endpoint MNC Wallet API.
//
//	Public:   POST /register, /login, /refresh, GET /health
//	Private:  POST /topup, /pay, /transfer
//	          GET  /transactions, /profile
//	          PUT  /profile
//
// Endpoint private dilindungi JWT middleware.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtSecret string) {
	RegisterHealth(r)

	r.POST("/register", h.Auth.Register)
	r.POST("/login", h.Auth.Login)
	r.POST("/refresh", h.Auth.Refresh)

	auth := r.Group("")
	auth.Use(middleware.JWTAuth(jwtSecret))
	{
		auth.POST("/topup", h.Wallet.TopUp)
		auth.POST("/pay", h.Wallet.Pay)
		auth.POST("/transfer", h.Transfer.Transfer)
		auth.GET("/transactions", h.Transaction.List)
		auth.GET("/profile", h.Profile.Get)
		auth.PUT("/profile", h.Profile.Update)
	}
}