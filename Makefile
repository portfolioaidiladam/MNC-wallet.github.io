# ---------------------------------------------------------------------------
# MNC Wallet API — Makefile
# Semua target membaca variable dari .env kalau ada.
# ---------------------------------------------------------------------------

ifneq (,$(wildcard .env))
	include .env
	export
endif

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= mncwallet
POSTGRES_PASSWORD ?= mncwallet
POSTGRES_DB ?= mncwallet
POSTGRES_SSLMODE ?= disable

MIGRATE_PATH := migrations
DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: help run worker test build migrate-up migrate-down migrate-new docker-up docker-down tidy fmt vet

help:
	@echo "Targets:"
	@echo "  make run           - run HTTP API (cmd/api)"
	@echo "  make worker        - run Asynq worker (cmd/worker)"
	@echo "  make test          - go test ./..."
	@echo "  make build         - compile both binaries into ./bin"
	@echo "  make migrate-up    - apply all pending migrations"
	@echo "  make migrate-down  - revert the last migration"
	@echo "  make migrate-new name=<desc> - scaffold new migration file pair"
	@echo "  make docker-up     - start postgres, redis, asynqmon"
	@echo "  make docker-down   - stop docker services"
	@echo "  make tidy          - go mod tidy"
	@echo "  make fmt           - gofmt -w ."
	@echo "  make vet           - go vet ./..."

run:
	go run ./cmd/api

worker:
	go run ./cmd/worker

test:
	go test ./...

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-new:
	@if [ -z "$(name)" ]; then echo "usage: make migrate-new name=<description>"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...