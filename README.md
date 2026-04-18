# MNC Wallet API

REST API e-wallet — technical test PT MNC Teknologi Nusantara (Fullstack Developer).
8 endpoint: register, login, refresh, top up, payment, transfer (async via Asynq),
transactions report, profile (GET + PUT).

## Tech Stack

- Go 1.22+, Gin, GORM v2, PostgreSQL 15, Redis 7, Asynq (+ Asynqmon)
- JWT v5 HS256 (access 15m / refresh 7d, refresh token disimpan sebagai SHA-256 hash)
- bcrypt cost 10
- golang-migrate untuk migrations
- testify untuk unit test

## Architecture

Clean architecture dependency flow: `handler -> service -> repository -> model`.
- **Handler**: parse request + map error ke HTTP status.
- **Service**: business logic, transactions, orkestrasi Asynq.
- **Repository**: pure DB access, row-level locking via `SELECT ... FOR UPDATE`.
- **Worker**: consume Asynq task `transfer:credit` untuk credit receiver.

## Requirements

- Go 1.22+
- Docker + Docker Compose (Postgres + Redis + Asynqmon)
- `migrate` CLI: `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- GNU Make (opsional)

## Setup Lokal

```bash
cp .env.example .env
# edit .env — wajib set JWT_SECRET

make docker-up       # postgres + redis + asynqmon
make migrate-up      # apply migrations
psql "$DB_URL" -f example_data.sql   # opsional, seed sample users (PIN 123456)

# Dua terminal:
make run             # HTTP API   :8000
make worker          # Asynq worker
```

Dashboard Asynqmon: http://localhost:8080 (tasks queue).

## Struktur Project

```
mnc-wallet/
├── cmd/api/              # HTTP server
├── cmd/worker/           # Asynq worker
├── internal/
│   ├── config/           # env loader + gorm open
│   ├── handler/          # HTTP handlers + router
│   ├── service/          # business logic (auth, wallet, transfer, ...)
│   ├── repository/       # GORM repo dengan FOR UPDATE locking
│   ├── model/            # domain + DTO
│   ├── middleware/       # JWT auth
│   ├── worker/           # Asynq task definitions + handlers
│   └── util/             # hash, JWT, response, validator, timestamp
├── migrations/           # golang-migrate SQL files
├── example_data.sql      # seed 3 user + sample transaksi
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Endpoint

| Method | Path           | Auth | Status  | Deskripsi                             |
|--------|----------------|------|---------|---------------------------------------|
| POST   | /register      | -    | 201     | Registrasi user baru + wallet kosong  |
| POST   | /login         | -    | 200     | Login, return access + refresh        |
| POST   | /refresh       | -    | 200     | Tukar refresh dengan pasangan baru    |
| POST   | /topup         | JWT  | 201     | Top up saldo (sync)                   |
| POST   | /pay           | JWT  | 201     | Payment merchant (sync)               |
| POST   | /transfer      | JWT  | 201     | Transfer antar-user (status PENDING)  |
| GET    | /transactions  | JWT  | 200     | Riwayat transaksi (LIMIT 100)         |
| GET    | /profile       | JWT  | 200     | Get profile + balance                 |
| PUT    | /profile       | JWT  | 200     | Update first_name/last_name/address   |

Response envelope:
```json
{ "status": "SUCCESS", "result": { ... } }
```
Error:
```json
{ "message": "human readable error" }
```

### Catatan desain

- **Transfer PENDING**: `POST /transfer` mendebit sender synchronously dalam DB
  transaction (row lock), lalu enqueue task Asynq `transfer:credit` untuk credit
  receiver. Response langsung dikirim dengan `status: "PENDING"` — field
  `balance_before` & `balance_after` sudah berisi saldo sender pasca-debit.
  Worker melakukan credit + menandai transaksi sender `SUCCESS` + insert
  `TRANSFER_IN` untuk receiver. Retry 3x exponential backoff; kalau tetap
  gagal, sender di-refund (compensating transaction) dan status diset `FAILED`.
- **Self-transfer ditolak** dengan 400 `cannot transfer to yourself`.
- **Phone**: prefix `08`, panjang 10–13 digit numerik.
- **PIN**: tepat 6 digit numerik, disimpan sebagai bcrypt hash (cost 10).
- **Refresh token rotation**: setiap `/refresh` revoke token lama + terbit
  sepasang token baru. Token server-side disimpan sebagai SHA-256 hex.
- **Transactions list**: sort DESC by `created_at`, LIMIT 100 (tidak paginate
  untuk scope test). `TransferID` di sisi `TRANSFER_IN` berisi `reference_id`
  yang sama dengan sisi `TRANSFER_OUT` — klien bisa match debit-credit.
- **Validasi amount**: hanya `amount > 0`. Production seharusnya punya min
  topup / daily limit — di luar scope test.

## Contoh curl

```bash
# Register
curl -sX POST localhost:8000/register \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"Aidil","last_name":"Adam","phone_number":"081234567890","address":"Jakarta","pin":"123456"}'

# Login
TOKENS=$(curl -sX POST localhost:8000/login \
  -H 'Content-Type: application/json' \
  -d '{"phone_number":"081234567890","pin":"123456"}')
ACCESS=$(echo "$TOKENS"  | jq -r .result.access_token)
REFRESH=$(echo "$TOKENS" | jq -r .result.refresh_token)

# Refresh
curl -sX POST localhost:8000/refresh \
  -H 'Content-Type: application/json' \
  -d "{\"refresh_token\":\"$REFRESH\"}"

# Top up
curl -sX POST localhost:8000/topup \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -d '{"amount":500000}'

# Payment
curl -sX POST localhost:8000/pay \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -d '{"amount":20000,"remarks":"beli kopi"}'

# Transfer (async)
curl -sX POST localhost:8000/transfer \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -d '{"target_user":"22222222-2222-2222-2222-222222222222","amount":100000,"remarks":"bayar hutang"}'

# Transactions
curl -s localhost:8000/transactions -H "Authorization: Bearer $ACCESS"

# Profile
curl -s localhost:8000/profile -H "Authorization: Bearer $ACCESS"
curl -sX PUT localhost:8000/profile \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"Aidil","last_name":"Adam","address":"Bandung"}'
```

## Sample Data

`example_data.sql` seed 3 user (PIN `123456`) + sample transaksi:

| User  | UserID                                   | Phone        | Balance |
|-------|------------------------------------------|--------------|---------|
| Aidil | 11111111-1111-1111-1111-111111111111     | 081234567890 | 380000  |
| Budi  | 22222222-2222-2222-2222-222222222222     | 081298765432 | 250000  |
| Citra | 33333333-3333-3333-3333-333333333333     | 081355512345 | 100000  |

Load:
```bash
psql postgres://mncwallet:mncwallet@localhost:5432/mncwallet?sslmode=disable \
  -f example_data.sql
```

## Testing

```bash
make test
```

Unit test coverage: util (jwt, validator, hash), service (auth, transaction list).
Happy-path `Register` dan semua operasi balance (topup/pay/transfer) butuh DB
asli karena pakai `SELECT ... FOR UPDATE` — verify via smoke test / curl di atas.

## Lisensi

Internal — technical test submission.