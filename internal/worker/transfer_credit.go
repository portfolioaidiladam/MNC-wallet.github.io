package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aidiladam/mnc-wallet/internal/model"
	"github.com/aidiladam/mnc-wallet/internal/repository"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// TransferCreditHandler memproses task TaskTypeTransferCredit: credit saldo
// receiver dan update TRANSFER_OUT sender jadi SUCCESS, plus insert TRANSFER_IN
// untuk receiver.
type TransferCreditHandler struct {
	db           *gorm.DB
	wallets      repository.WalletRepository
	transactions repository.TransactionRepository
}

// NewTransferCreditHandler membangun handler yang diserahkan ke asynq mux.
func NewTransferCreditHandler(
	db *gorm.DB,
	wallets repository.WalletRepository,
	transactions repository.TransactionRepository,
) *TransferCreditHandler {
	return &TransferCreditHandler{db: db, wallets: wallets, transactions: transactions}
}

// Handle mengimplementasikan asynq.HandlerFunc.
//
// Retry policy: asynq auto retry kalau kita return error. Di attempt ke-N
// (max retries tercapai), kita panggil compensating transaction: refund sender
// dan tandai transaksi FAILED. Asynq lalu anggap task sukses — supaya tidak
// di-retry lagi.
func (h *TransferCreditHandler) Handle(ctx context.Context, task *asynq.Task) error {
	var p TransferCreditPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		// Payload rusak — tidak ada gunanya retry.
		return fmt.Errorf("unmarshal payload: %w: %w", err, asynq.SkipRetry)
	}

	attempt, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)

	if err := h.processCredit(ctx, &p); err != nil {
		log.Printf("worker: transfer credit attempt %d/%d failed ref=%s err=%v",
			attempt+1, maxRetry+1, p.ReferenceID, err)

		if attempt >= maxRetry {
			// Final failure: refund sender + mark transaksi FAILED,
			// lalu return nil supaya asynq tidak retry lagi.
			if cerr := h.compensate(ctx, &p, err.Error()); cerr != nil {
				log.Printf("worker: compensation failed ref=%s err=%v", p.ReferenceID, cerr)
				return cerr
			}
			log.Printf("worker: transfer credit compensated (refund sender) ref=%s", p.ReferenceID)
			return nil
		}
		return err
	}
	log.Printf("worker: transfer credit OK ref=%s sender=%s receiver=%s amount=%d",
		p.ReferenceID, p.SenderID, p.ReceiverID, p.Amount)
	return nil
}

// processCredit menjalankan credit receiver + mark sender SUCCESS + insert
// TRANSFER_IN dalam satu DB transaction. Idempotent terhadap retry: kalau
// TRANSFER_IN dengan reference_id yang sama sudah ada, skip.
func (h *TransferCreditHandler) processCredit(ctx context.Context, p *TransferCreditPayload) error {
	return repository.RunInTx(ctx, h.db, func(gtx *gorm.DB) error {
		// Idempotency guard: cek apakah TRANSFER_IN untuk reference ini sudah ada.
		var existing int64
		if err := gtx.WithContext(ctx).Model(&model.Transaction{}).
			Where("reference_id = ? AND type = ?", p.ReferenceID, model.TxTypeTransferIn).
			Count(&existing).Error; err != nil {
			return fmt.Errorf("check idempotency: %w", err)
		}
		if existing > 0 {
			return nil
		}

		// Lock receiver wallet + credit.
		rw, err := h.wallets.LockByUserID(ctx, gtx, p.ReceiverID)
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("receiver wallet not found")
		}
		if err != nil {
			return err
		}
		newBalance := rw.Balance + p.Amount
		if err := h.wallets.UpdateBalance(ctx, gtx, rw.ID, newBalance); err != nil {
			return err
		}

		// Insert TRANSFER_IN.
		txIn := &model.Transaction{
			ID:                 model.NewID(),
			UserID:             p.ReceiverID,
			Type:               model.TxTypeTransferIn,
			Status:             model.TxStatusSuccess,
			Amount:             p.Amount,
			BalanceBefore:      rw.Balance,
			BalanceAfter:       newBalance,
			Remarks:            p.Remarks,
			ReferenceID:        &p.ReferenceID,
			CounterpartyUserID: &p.SenderID,
			CreatedAt:          time.Now(),
		}
		if err := h.transactions.Create(ctx, gtx, txIn); err != nil {
			return err
		}

		// Flip TRANSFER_OUT sender jadi SUCCESS.
		if err := gtx.WithContext(ctx).Model(&model.Transaction{}).
			Where("id = ?", p.ReferenceID).
			Update("status", model.TxStatusSuccess).Error; err != nil {
			return fmt.Errorf("update sender transaction: %w", err)
		}
		return nil
	})
}

// compensate me-refund sender dan menandai transaksi debit FAILED.
// Dipanggil hanya di final attempt kalau processCredit gagal terus.
func (h *TransferCreditHandler) compensate(ctx context.Context, p *TransferCreditPayload, reason string) error {
	return repository.RunInTx(ctx, h.db, func(gtx *gorm.DB) error {
		// Idempotency: kalau sudah FAILED, skip.
		var current model.Transaction
		if err := gtx.WithContext(ctx).Where("id = ?", p.ReferenceID).First(&current).Error; err != nil {
			return fmt.Errorf("load sender transaction: %w", err)
		}
		if current.Status == model.TxStatusFailed {
			return nil
		}
		if current.Status == model.TxStatusSuccess {
			// Sudah berhasil di attempt lain (seharusnya tidak terjadi di jalur ini).
			return nil
		}

		sw, err := h.wallets.LockByUserID(ctx, gtx, p.SenderID)
		if err != nil {
			return err
		}
		newBalance := sw.Balance + p.Amount
		if err := h.wallets.UpdateBalance(ctx, gtx, sw.ID, newBalance); err != nil {
			return err
		}
		newRemarks := current.Remarks
		if newRemarks != "" {
			newRemarks += " | "
		}
		newRemarks += "refund: " + reason
		if len(newRemarks) > 255 {
			newRemarks = newRemarks[:255]
		}
		return gtx.WithContext(ctx).Model(&model.Transaction{}).
			Where("id = ?", p.ReferenceID).
			Updates(map[string]any{
				"status":  model.TxStatusFailed,
				"remarks": newRemarks,
			}).Error
	})
}