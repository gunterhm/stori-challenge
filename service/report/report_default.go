package report

import (
	"awesomeProject/model"
	"awesomeProject/repository/accountrepo"
	"context"
	"go.uber.org/zap"
)

// DefaultService is an implementation of IReport
type DefaultService struct {
	log         *zap.SugaredLogger
	accountRepo accountrepo.IRepository
}

// NewDefaultService creates a new DefaultService
func NewDefaultService(logger *zap.SugaredLogger, accRepo accountrepo.IRepository) *DefaultService {
	return &DefaultService{
		log:         logger,
		accountRepo: accRepo,
	}
}

func (s DefaultService) GetSummaryEmailInfo(ctx context.Context, accountID string) (*model.SummaryEmail, error) {
	account, err := s.accountRepo.FindAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	averageCredit, averageDebit, err := s.accountRepo.FindMonthlyAverageByTxnType(ctx, accountID)
	if err != nil {
		return nil, err
	}

	arrTxnCountPerMonth, err := s.accountRepo.FindTotalTransactionsPerMonth(ctx, accountID)
	if err != nil {
		return nil, err
	}

	totalBalance, err := s.accountRepo.FindTotalBalance(ctx, accountID)
	if err != nil {
		return nil, err
	}

	summaryEmail := model.SummaryEmail{
		AccountNumber:       account.AccountID,
		CustomerName:        account.Name,
		Email:               account.Email,
		TotalBalance:        totalBalance,
		AverageCreditAmount: averageCredit,
		AverageDebitAmount:  averageDebit,
		ArrTxnCountPerMonth: arrTxnCountPerMonth,
	}

	return &summaryEmail, nil
}
