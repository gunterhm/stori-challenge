// Package account contains repository functions
package account

import (
	"awesomeProject/model"
	"context"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// DatabaseRepository  repository implementation
type DatabaseRepository struct {
	log *zap.SugaredLogger
	db  *bun.DB
}

// NewDatabaseRepository Creates a new repository
func NewDatabaseRepository(logger *zap.SugaredLogger, conn *bun.DB) *DatabaseRepository {
	return &DatabaseRepository{
		log: logger,
		db:  conn,
	}
}

// FindAccount Finds account
func (r *DatabaseRepository) FindAccount(ctx context.Context, accountID string) (*model.Account, error) {
	account := new(model.Account)
	account.AccountID = accountID

	err := r.db.NewSelect().
		Model(account).
		WherePK().
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// FindMonthlyAverageByTxnType Finds monthly average by transaction type
func (r *DatabaseRepository) FindMonthlyAverageByTxnType(ctx context.Context, accountID string) (float64, float64, error) {

	subQuery := r.db.NewSelect().
		Model((*model.AccountTransaction)(nil)).
		ColumnExpr("MONTH(date) AS txn_month").
		ColumnExpr("SUM(amount_credit) AS sum_txn_amount_credit").
		ColumnExpr("SUM(amount_debit) AS sum_txn_amount_debit").
		Where("account_id = ?", accountID).
		Where("YEAR(date) = YEAR(CURDATE())").
		GroupExpr("MONTH(date)")

	var averageCredit float64
	var averageDebit float64
	err := r.db.NewSelect().
		ColumnExpr("AVG(sum_txn_amount_credit)").
		ColumnExpr("AVG(sum_txn_amount_debit)").
		TableExpr("(?) AS monthly_totals", subQuery).
		Scan(ctx, &averageCredit, &averageDebit)

	if err != nil {
		return 0, 0, err
	}

	return averageCredit, averageDebit, nil
}
