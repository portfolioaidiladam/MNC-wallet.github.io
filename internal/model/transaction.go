package model

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType adalah ENUM transaction_type di Postgres.
type TransactionType string

// TransactionStatus adalah ENUM transaction_status di Postgres.
type TransactionStatus string

const (
	TxTypeTopup       TransactionType = "TOPUP"
	TxTypePayment     TransactionType = "PAYMENT"
	TxTypeTransferOut TransactionType = "TRANSFER_OUT"
	TxTypeTransferIn  TransactionType = "TRANSFER_IN"

	TxStatusPending TransactionStatus = "PENDING"
	TxStatusSuccess TransactionStatus = "SUCCESS"
	TxStatusFailed  TransactionStatus = "FAILED"
)

// Transaction adalah satu baris riwayat di tabel transactions.
//
// ReferenceID menghubungkan dua sisi TRANSFER (TRANSFER_OUT dan TRANSFER_IN
// yang berpasangan share satu reference_id). Null untuk TOPUP/PAYMENT.
//
// CounterpartyUserID berisi user_id lawan untuk transfer; null untuk
// TOPUP/PAYMENT.
type Transaction struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey;column:id"`
	UserID             uuid.UUID         `gorm:"type:uuid;not null;column:user_id"`
	Type               TransactionType   `gorm:"type:transaction_type;not null;column:type"`
	Status             TransactionStatus `gorm:"type:transaction_status;not null;default:'PENDING';column:status"`
	Amount             int64             `gorm:"not null;column:amount"`
	BalanceBefore      int64             `gorm:"not null;column:balance_before"`
	BalanceAfter       int64             `gorm:"not null;column:balance_after"`
	Remarks            string            `gorm:"size:255;not null;default:'';column:remarks"`
	ReferenceID        *uuid.UUID        `gorm:"type:uuid;column:reference_id"`
	CounterpartyUserID *uuid.UUID        `gorm:"type:uuid;column:counterparty_user_id"`
	CreatedAt          time.Time         `gorm:"not null;column:created_at"`
}

// TableName mengoverride nama tabel default GORM.
func (Transaction) TableName() string { return "transactions" }