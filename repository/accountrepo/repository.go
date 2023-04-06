package accountrepo

import (
	"awesomeProject/model"
	"context"
)

// IRepository interface
type IRepository interface {
	FindAccount(ctx context.Context, accountID string) (*model.Account, error)
	FindMonthlyAverageByTxnType(ctx context.Context, accountID string) (float64, float64, error)
	FindTotalTransactionsPerMonth(ctx context.Context, accountID string) ([]model.TxnCountPerMonth, error)
}
