package service

import (
	"context"
	"fmt"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/google/uuid"
)

// TransactionListLimit adalah batas maksimum record yang dikembalikan /transactions.
// Hardcode (tidak paginate) untuk scope technical test.
const TransactionListLimit = 100

// TransactionService menyajikan riwayat transaksi user.
type TransactionService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.TransactionListItem, error)
}

type transactionService struct {
	transactions repository.TransactionRepository
}

// NewTransactionService membangun TransactionService.
func NewTransactionService(transactions repository.TransactionRepository) TransactionService {
	return &transactionService{transactions: transactions}
}

// List mengembalikan transaksi user (TOPUP, PAYMENT, TRANSFER_OUT, TRANSFER_IN),
// di-map ke DTO TransactionListItem dengan field ID per-tipe.
func (s *transactionService) List(ctx context.Context, userID uuid.UUID) ([]model.TransactionListItem, error) {
	rows, err := s.transactions.ListByUserID(ctx, userID, TransactionListLimit)
	if err != nil {
		return nil, err
	}
	out := make([]model.TransactionListItem, 0, len(rows))
	for i := range rows {
		item, err := toListItem(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// toListItem mengkonversi Transaction ke TransactionListItem; salah satu dari
// TopUpID/PaymentID/TransferID terisi sesuai Type.
func toListItem(t *model.Transaction) (model.TransactionListItem, error) {
	item := model.TransactionListItem{
		Status:          string(t.Status),
		UserID:          t.UserID,
		TransactionType: string(t.Type),
		Amount:          t.Amount,
		Remarks:         t.Remarks,
		BalanceBefore:   t.BalanceBefore,
		BalanceAfter:    t.BalanceAfter,
		CreatedAt:       util.FormatJakarta(t.CreatedAt),
	}
	id := t.ID
	switch t.Type {
	case model.TxTypeTopup:
		item.TopUpID = &id
	case model.TxTypePayment:
		item.PaymentID = &id
	case model.TxTypeTransferOut, model.TxTypeTransferIn:
		// ReferenceID = id debit-side. Dipakai di kedua sisi supaya
		// client bisa trace pair debit-credit.
		ref := t.ID
		if t.ReferenceID != nil {
			ref = *t.ReferenceID
		}
		item.TransferID = &ref
	default:
		return item, fmt.Errorf("unknown transaction type %q", t.Type)
	}
	return item, nil
}