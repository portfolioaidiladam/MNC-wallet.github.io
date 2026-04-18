// Package worker berisi Asynq task handlers (mis. transfer credit receiver).
package worker

// TaskTypeTransferCredit adalah type name Asynq task untuk credit balance
// receiver pada transfer antar-user.
const TaskTypeTransferCredit = "transfer:credit"