// Package accountrepo contains repository functions for accounts in database
package accountrepo

import (
	"awesomeProject/model"
	"context"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// DatabaseRepository is an implementation of IRepository
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

// FindMonthlyAverageByTxnType Finds monthly credit and debit average
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

// FindTotalTransactionsPerMonth Finds total transactions per month
func (r *DatabaseRepository) FindTotalTransactionsPerMonth(ctx context.Context, accountID string) ([]model.TxnCountPerMonth, error) {

	var txnCountPerMonth []model.TxnCountPerMonth

	err := r.db.NewSelect().
		Model((*model.AccountTransaction)(nil)).
		ColumnExpr("MONTH(date) AS txn_month").
		ColumnExpr("COUNT(*) AS txn_count").
		Where("account_id = ?", accountID).
		Where("YEAR(date) = YEAR(CURDATE())").
		GroupExpr("MONTH(date)").
		OrderExpr("txn_month ASC").
		Scan(ctx, &txnCountPerMonth)

	if err != nil {
		return nil, err
	}

	return txnCountPerMonth, nil
}

func (r *DatabaseRepository) FindTotalBalance(ctx context.Context, accountID string) (float64, error) {

	var totalBalance float64

	err := r.db.NewSelect().
		Model((*model.AccountTransaction)(nil)).
		ColumnExpr("SUM(amount_credit) + SUM(amount_debit) AS total_balance").
		Where("account_id = ?", accountID).
		Where("YEAR(date) = YEAR(CURDATE())").
		Scan(ctx, &totalBalance)

	if err != nil {
		return 0, err
	}

	return totalBalance, nil
}

func (r *DatabaseRepository) SaveAccountTransaction(ctx context.Context, transaction *model.AccountTransaction) error {
	_, err := r.db.NewInsert().
		On("DUPLICATE KEY UPDATE").
		Model(transaction).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
