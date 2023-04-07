package main

import (
	"awesomeProject/config"
	"awesomeProject/repository/accountrepo"
	"awesomeProject/service/mail"
	"awesomeProject/service/report"
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
	reportSvc := report.NewDefaultService(logger, repoAccount)
	mailSvc := mail.NewDefaultService(logger, &configs.Smtp)

	txnProcessorSvc := txnprocessor.NewDefaultService(logger,
		repoAccount,
		reportSvc,
		mailSvc,
		configs.TxnProcessing.FileNameRegExp,
		configs.TxnProcessing.IncomingDir, configs.TxnProcessing.ArchiveDir)

	// Go Cron
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(configs.Schedule.Seconds).Seconds().
		Do(func() {
			err = txnProcessorSvc.StartProcess(context.Background())
			if err != nil {
				logger.Errorf("Error while calling StartProcess: %v", err)
			}
		})

	scheduler.StartBlocking()
}
