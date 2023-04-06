package account

import (
	"awesomeProject/model"
	"context"
)

// IRepository interface
type IRepository interface {
	FindAccount(ctx context.Context, accountID string) (*model.Account, error)
	FindMonthlyAverageByTxnType(ctx context.Context, accountID string, txnType string) (float64, error)
}
