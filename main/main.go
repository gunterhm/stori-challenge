package main

import (
	"awesomeProject/config"
	"awesomeProject/repository/accountrepo"
	"awesomeProject/service/txnprocessor"
	"context"
	"database/sql"
	"github.com/go-co-op/gocron"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
	"time"
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
	/* TODO GUNTER TO UNCOMMENT
	reportSvc := report.NewDefaultService(logger, repoAccount)
	mailSvc := mail.NewDefaultService(logger, &configs.Smtp)
	*/
	txnProcessorSvc := txnprocessor.NewDefaultService(logger,
		repoAccount,
		configs.TxnProcessing.FileNameRegExp,
		configs.TxnProcessing.IncomingDir, configs.TxnProcessing.ArchiveDir)

	// Go Cron
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(5).Seconds().
		Do(func() { logger.Infof("TICK!!!!") })

	err = txnProcessorSvc.StartProcess(context.Background())
	if err != nil {
		logger.Errorf("Error while calling StartProcess: %v", err)
	}

	/*  TODO GUNTER TO UNCOMMENT
	summaryEmailInfo, err := reportSvc.GetSummaryEmailInfo(context.Background(), "GHM54789345")
	if err != nil {
		return
	}

	logger.Infof("Summary Email info %v", summaryEmailInfo)

	err = mailSvc.SendSummaryMail(summaryEmailInfo)
	if err != nil {
		return
	}*/

	scheduler.StartBlocking()
}
