package entity

import (
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
)

type Merchant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transfer struct {
	ID             string                  `json:"id"`
	MerchantID     int64                   `json:"merchant_id"`
	IdempotencyKey string                  `json:"idempotency_key"`
	PartnerID      string                  `json:"partner_id"`
	BankCode       string                  `json:"bank_code"`
	AccountNumber  string                  `json:"account_number"`
	Amount         int64                   `json:"amount"`
	Status         constant.TransferStatus `json:"status"`
	RetryCount     int                     `json:"retry_count"`
	FailureReason  *string                 `json:"failure_reason,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type BalanceMutation struct {
	ID            int64                 `json:"id"`
	MerchantID    int64                 `json:"merchant_id"`
	TransferID    string                `json:"transfer_id"`
	Type          constant.MutationType `json:"type"`
	Amount        int64                 `json:"amount"`
	BalanceBefore int64                 `json:"balance_before"`
	BalanceAfter  int64                 `json:"balance_after"`
	CreatedAt     time.Time             `json:"created_at"`
}

type CreateTransferRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	BankCode       string `json:"bank_code"`
	AccountNumber  string `json:"account_number"`
	Amount         int64  `json:"amount"`
}

type CreateTransferResponse struct {
	TransferID string                  `json:"transfer_id"`
	Status     constant.TransferStatus `json:"status"`
	Amount     int64                   `json:"amount"`
}

type TransferMessage struct {
	TransferID string `json:"transfer_id"`
}
