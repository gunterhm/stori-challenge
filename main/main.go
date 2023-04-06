package main

import (
	"awesomeProject/config"
	"awesomeProject/model"
	"awesomeProject/repository/account"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func main() {
	// Logger
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	// Configs
	configs := config.LoadConfig(logger)

	// Database
	sqlDb, err := sql.Open("mysql", configs.Database.Dsn)
	if err != nil {
		logger.Errorf("Connection to database failed: %v", err)
	}
	defer sqlDb.Close()

	logger.Infof("Connected to Database!")

	db := bun.NewDB(sqlDb, mysqldialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Repositories
	repoAccount := account.NewDatabaseRepository(logger, db)

	account, err := repoAccount.FindAccount(context.Background(), "GHM54789345")
	if err != nil {
		logger.Errorf("Error in FindAccount: %v", err)
	}

	logger.Infof("Account: %v", account)

	avgCredit, avgDebit, err := repoAccount.FindMonthlyAverageByTxnType(context.Background(), "GHM54789345")
	if err != nil {
		logger.Errorf("Error in FindMonthlyAverageByTxnType: %v", err)
	}

	logger.Infof("Monthly Average (CREDIT): %f  Monthly Average (DEBIT): %f", avgCredit, avgDebit)

	// Services
	//mailSvc := mail.NewDefaultService(logger, &configs.Smtp)  TODO GUNTER TO UNCOMMENT

	var arrTxnCountPerMonth = []model.TxnCountPerMonth{
		{
			Month:    "January",
			TxnCount: 123,
		},
		{
			Month:    "February",
			TxnCount: 456,
		},
		{
			Month:    "March",
			TxnCount: 789,
		},
	}

	var summaryEmailData = model.SummaryEmail{
		AccountNumber:       "GHM09238458",
		TotalBalance:        540.8,
		AverageCreditAmount: 105.5,
		AverageDebitAmount:  -56.3,
		ArrTxnCountPerMonth: arrTxnCountPerMonth,
	}

	logger.Infof("Summary Email Data: %v", summaryEmailData)

	/*err = mailSvc.SendMail("gunterhm@gmail.com", summaryEmailData)  TODO GUNTER TO UNCOMENT
	if err != nil {
		return
	}*/
}
