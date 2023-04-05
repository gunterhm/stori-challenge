package main

import (
	"awesomeProject/config"
	"awesomeProject/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
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
		logger.Error("Connection to database failed: %v", err)
	}
	defer sqlDb.Close()

	logger.Info("Connected to Database!")

	_ = bun.NewDB(sqlDb, mysqldialect.New())

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
