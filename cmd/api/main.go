// Command api menjalankan HTTP server MNC Wallet.
//
// Responsibility: load config, open Postgres + Redis (Asynq client), wire
// layers (repository -> service -> handler), dan expose HTTP routes.
package main

import (
	"log"

	"github.com/aidiladam/mnc-wallet/internal/config"
	"github.com/aidiladam/mnc-wallet/internal/handler"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := config.OpenDatabase(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer func() { _ = asynqClient.Close() }()

	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	rtRepo := repository.NewRefreshTokenRepository(db)

	authSvc := service.NewAuthService(db, userRepo, walletRepo, rtRepo, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	walletSvc := service.NewWalletService(db, walletRepo, txRepo)
	transferSvc := service.NewTransferService(db, userRepo, walletRepo, txRepo, asynqClient)
	txSvc := service.NewTransactionService(txRepo)
	profileSvc := service.NewProfileService(userRepo, walletRepo)

	h := &handler.Handlers{
		Auth:        handler.NewAuthHandler(authSvc),
		Wallet:      handler.NewWalletHandler(walletSvc),
		Transfer:    handler.NewTransferHandler(transferSvc),
		Transaction: handler.NewTransactionHandler(txSvc),
		Profile:     handler.NewProfileHandler(profileSvc),
	}

	r := gin.Default()
	handler.RegisterRoutes(r, h, cfg.JWTSecret)

	log.Printf("api listening on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}