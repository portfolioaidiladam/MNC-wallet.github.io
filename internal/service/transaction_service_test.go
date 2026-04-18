package service

import (
	"context"
	"testing"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionService_List_MapsPerTypeID(t *testing.T) {
	userID := uuid.New()
	counterparty := uuid.New()

	topupID := uuid.New()
	paymentID := uuid.New()
	transferOutID := uuid.New()
	transferInID := uuid.New()
	ref := transferOutID

	txRepo := &mockTxRepo{}
	rows := []model.Transaction{
		{
			ID: topupID, UserID: userID, Type: model.TxTypeTopup, Status: model.TxStatusSuccess,
			Amount: 100_000, BalanceBefore: 0, BalanceAfter: 100_000,
			CreatedAt: time.Now(),
		},
		{
			ID: paymentID, UserID: userID, Type: model.TxTypePayment, Status: model.TxStatusSuccess,
			Amount: 20_000, BalanceBefore: 100_000, BalanceAfter: 80_000,
			Remarks: "kopi", CreatedAt: time.Now(),
		},
		{
			ID: transferOutID, UserID: userID, Type: model.TxTypeTransferOut, Status: model.TxStatusSuccess,
			Amount: 50_000, BalanceBefore: 80_000, BalanceAfter: 30_000,
			ReferenceID: &ref, CounterpartyUserID: &counterparty, CreatedAt: time.Now(),
		},
		{
			ID: transferInID, UserID: userID, Type: model.TxTypeTransferIn, Status: model.TxStatusSuccess,
			Amount: 75_000, BalanceBefore: 30_000, BalanceAfter: 105_000,
			ReferenceID: &ref, CounterpartyUserID: &counterparty, CreatedAt: time.Now(),
		},
	}
	txRepo.On("ListByUserID", mock.Anything, userID, TransactionListLimit).Return(rows, nil)

	svc := NewTransactionService(txRepo)
	out, err := svc.List(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, out, 4)

	// TOPUP: hanya top_up_id terisi
	assert.NotNil(t, out[0].TopUpID)
	assert.Equal(t, topupID, *out[0].TopUpID)
	assert.Nil(t, out[0].PaymentID)
	assert.Nil(t, out[0].TransferID)

	// PAYMENT
	assert.NotNil(t, out[1].PaymentID)
	assert.Equal(t, paymentID, *out[1].PaymentID)
	assert.Nil(t, out[1].TopUpID)
	assert.Nil(t, out[1].TransferID)

	// TRANSFER_OUT: TransferID = reference_id (= id itu sendiri di debit side)
	assert.NotNil(t, out[2].TransferID)
	assert.Equal(t, transferOutID, *out[2].TransferID)
	assert.Nil(t, out[2].TopUpID)
	assert.Nil(t, out[2].PaymentID)

	// TRANSFER_IN: TransferID = reference_id (sama dengan sisi OUT)
	assert.NotNil(t, out[3].TransferID)
	assert.Equal(t, transferOutID, *out[3].TransferID)

	txRepo.AssertExpectations(t)
}

func TestTransactionService_List_Empty(t *testing.T) {
	userID := uuid.New()
	txRepo := &mockTxRepo{}
	txRepo.On("ListByUserID", mock.Anything, userID, TransactionListLimit).
		Return([]model.Transaction{}, nil)

	svc := NewTransactionService(txRepo)
	out, err := svc.List(context.Background(), userID)
	require.NoError(t, err)
	assert.Empty(t, out)
}