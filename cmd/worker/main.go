// Command worker menjalankan Asynq worker MNC Wallet.
//
// Responsibility: load config, open Postgres + Redis, daftarkan task handler,
// dan blocking run asynq.Server.
package main

import (
	"log"

	"github.com/aidiladam/mnc-wallet/internal/config"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/worker"
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

	walletRepo := repository.NewWalletRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	creditHandler := worker.NewTransferCreditHandler(db, walletRepo, txRepo)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
		},
		asynq.Config{
			Concurrency: cfg.AsynqConcurrency,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TaskTypeTransferCredit, creditHandler.Handle)

	log.Printf("asynq worker starting (concurrency=%d, broker=%s)", cfg.AsynqConcurrency, cfg.RedisAddr)
	if err := srv.Run(mux); err != nil {
		log.Fatalf("asynq: %v", err)
	}
}