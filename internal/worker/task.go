package worker

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// TransferCreditPayload adalah payload JSON task TaskTypeTransferCredit.
//
// Semua field id yang memang wajib di-populate di enqueue time. ReferenceID
// sama dengan id TRANSFER_OUT supaya worker bisa match side debit-nya.
type TransferCreditPayload struct {
	ReferenceID uuid.UUID `json:"reference_id"`
	SenderID    uuid.UUID `json:"sender_id"`
	ReceiverID  uuid.UUID `json:"receiver_id"`
	Amount      int64     `json:"amount"`
	Remarks     string    `json:"remarks"`
}

// NewTransferCreditTask membangun Asynq task dengan payload JSON.
// Retry dan queue di-set di enqueuer (service layer) supaya terpisah
// dari definisi task itu sendiri.
func NewTransferCreditTask(p TransferCreditPayload) (*asynq.Task, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshal transfer credit payload: %w", err)
	}
	return asynq.NewTask(TaskTypeTransferCredit, b), nil
}