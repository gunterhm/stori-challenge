package txnprocessor

import (
	"awesomeProject/model"
	"awesomeProject/repository/accountrepo"
	"awesomeProject/service/mail"
	"awesomeProject/service/report"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

// DefaultService is an implementation of ITxnProcessorService
type DefaultService struct {
	log            *zap.SugaredLogger
	accountRepo    accountrepo.IRepository
	reportSvc      report.IReport
	mailSvc        mail.IMailService
	fileNameRegExp *regexp.Regexp
	incomingDir    string
	archiveDir     string
}

// NewDefaultService creates a new DefaultService
func NewDefaultService(logger *zap.SugaredLogger,
	accRepo accountrepo.IRepository, repSvc report.IReport, mailSvc mail.IMailService,
	strFileNameRegex string, incDir string, archDir string) *DefaultService {
	rexExp, _ := regexp.Compile(strFileNameRegex)
	return &DefaultService{
		log:            logger,
		accountRepo:    accRepo,
		reportSvc:      repSvc,
		mailSvc:        mailSvc,
		fileNameRegExp: rexExp,
		incomingDir:    incDir,
		archiveDir:     archDir,
	}
}

func (s DefaultService) StartProcess(ctx context.Context) error {
	nextFile, err := s.NextTxnFile()
	if err != nil {
		return err
	}

	if nextFile == nil {
		s.log.Infof("No files to process.")
		return nil
	}
	s.log.Infof("Next File to process: %v", *nextFile)

	err = s.ProcessTxnFile(ctx, nextFile)
	if err != nil {
		return err
	}

	summaryEmailInfo, err := s.reportSvc.GetSummaryEmailInfo(context.Background(), "GHM54789345")
	if err != nil {
		return err
	}

	s.log.Infof("Summary Email info %v", summaryEmailInfo)

	err = s.mailSvc.SendSummaryMail(summaryEmailInfo)
	if err != nil {
		return err
	}

	return nil
}

func (s DefaultService) NextTxnFile() (*os.DirEntry, error) {
	files, err := os.ReadDir(s.incomingDir)
	if err != nil {
		return nil, err
	}

	var nextFile *os.DirEntry = nil
	for _, file := range files {
		if !file.IsDir() {
			if s.fileNameRegExp.MatchString(file.Name()) {
				nextFile = &file
				break
			}
		}
	}

	return nextFile, nil
}

func (s DefaultService) ProcessTxnFile(ctx context.Context, txnFile *os.DirEntry) error {
	// Rename txn file while being processed
	originalFilename := (*txnFile).Name()
	processingFilename := "PROCESSING_" + originalFilename
	s.log.Infof("Original Name: %s", originalFilename)
	err := os.Rename(s.incomingDir+string(os.PathSeparator)+originalFilename, s.incomingDir+string(os.PathSeparator)+processingFilename)
	if err != nil {
		return err
	}

	//Extract account ID from file name
	matches := s.fileNameRegExp.FindStringSubmatch(originalFilename)
	accountID := matches[1]
	s.log.Infof("Extracted Account ID: %s", accountID)

	// Open txn file
	f, err := os.Open(s.incomingDir + string(os.PathSeparator) + processingFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read csv values
	csvReader := csv.NewReader(f)
	lineNumber := 0
	for {
		lineNumber++
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Convert csv entry to AccountTransaction structure
		accountTransaction, err := mapCsvRecordToAccountTransaction(rec, accountID)
		if err != nil {
			s.log.Infof("Error at line %d: %v", lineNumber, err)
		} else {
			err := s.accountRepo.SaveAccountTransaction(ctx, accountTransaction)
			if err != nil {
				return err
			}
		}
		fmt.Printf("%+v\n", rec)
	}

	// Move txn file to archive
	err = os.Rename(s.incomingDir+string(os.PathSeparator)+processingFilename, s.archiveDir+string(os.PathSeparator)+originalFilename)
	if err != nil {
		return err
	}
	return nil
}

func mapCsvRecordToAccountTransaction(record []string, accountID string) (*model.AccountTransaction, error) {
	if len(record) != 3 {
		return nil, errors.New("invalid record size")
	}

	// Extract Txn Id from CSV record
	txnID, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, errors.New("txn ID: " + record[0] + " is not an integer")
	}

	// Extract Txn Amount from CSV record
	amount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, errors.New("txn amount: " + record[2] + " is not a number")
	}

	// Extract Txn Date from CSV record
	dateElements := strings.Split(record[1], "/")
	if len(dateElements) != 2 {
		return nil, errors.New("txn date: " + record[1] + " is not a date with format month/date")
	}
	if len(dateElements[0]) < 2 {
		dateElements[0] = "0" + dateElements[0]
	}
	if len(dateElements[1]) < 2 {
		dateElements[1] = "0" + dateElements[1]
	}
	dateString := strconv.Itoa(time.Now().Year()) + "-" + dateElements[0] + "-" + dateElements[1]
	txnDate, err := time.Parse(dateLayout, dateString)
	if err != nil {
		return nil, errors.New("txn date: " + record[1] + " is not a date with format month/date")
	}

	// Set Credit and Debit amounts
	var amountCredit, amountDebit float64
	if amount < 0 {
		amountCredit = 0
		amountDebit = amount
	} else {
		amountCredit = amount
		amountDebit = 0
	}

	// Fill in structure fields
	accountTxn := model.AccountTransaction{
		AccountID:     accountID,
		TransactionID: txnID,
		Date:          txnDate,
		AmountCredit:  amountCredit,
		AmountDebit:   amountDebit,
	}

	return &accountTxn, nil
}
