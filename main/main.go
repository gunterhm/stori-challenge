package main

import (
	"awesomeProject/config"
	"awesomeProject/repository/accountrepo"
	"awesomeProject/service/mail"
	"awesomeProject/service/report"
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
	repoAccount := accountrepo.NewDatabaseRepository(logger, db)

	// Services
	reportSvc := report.NewDefaultService(logger, repoAccount)
	mailSvc := mail.NewDefaultService(logger, &configs.Smtp)

	summaryEmailInfo, err := reportSvc.GetSummaryEmailInfo(context.Background(), "GHM54789345")
	if err != nil {
		return
	}

	logger.Infof("Summary Email info %v", summaryEmailInfo)

	err = mailSvc.SendMail("gunterhm@gmail.com", summaryEmailInfo)
	if err != nil {
		return
	}
}
