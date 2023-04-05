package model

import (
	"github.com/uptrace/bun"
)

// Account is the model for the account table
type Account struct {
	bun.BaseModel `bun:"table:account"`

	AccountID string `bun:"account_id,pk"`
	Name      string `bun:"name"`
	Email     string `bun:"email"`
}
