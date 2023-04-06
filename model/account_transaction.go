package model

import (
	"github.com/uptrace/bun"
	"time"
)

// AccountTransaction is the model for the account_transaction table
type AccountTransaction struct {
	bun.BaseModel `bun:"table:account_transaction"`

	AccountID     string    `bun:"account_id,pk"`
	TransactionID string    `bun:"transaction_id"`
	Date          time.Time `bun:"date"`
	AmountCredit  float64   `bun:"amount_credit"`
	AmountDebit   float64   `bun:"amount_debit"`
	Type          string    `bun:"type"`
	Account       *Account  `bun:"rel:belongs-to,join:account_id=account_id"`
}
