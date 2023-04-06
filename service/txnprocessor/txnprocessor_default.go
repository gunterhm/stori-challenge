package txnprocessor

import (
	"encoding/csv"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"regexp"
)

// DefaultService is an implementation of ITxnProcessorService
type DefaultService struct {
	log            *zap.SugaredLogger
	fileNameRegExp *regexp.Regexp
	incomingDir    string
	archiveDir     string
}

// NewDefaultService creates a new DefaultService
func NewDefaultService(logger *zap.SugaredLogger, strFileNameRegex string, incDir string, archDir string) *DefaultService {
	rexExp, _ := regexp.Compile(strFileNameRegex)
	return &DefaultService{
		log:            logger,
		fileNameRegExp: rexExp,
		incomingDir:    incDir,
		archiveDir:     archDir,
	}
}

func (s DefaultService) StartProcess() error {
	nextFile, err := s.NextTxnFile()
	if err != nil {
		return err
	}
	if nextFile == nil {
		s.log.Infof("No files to process.")
	}
	s.log.Infof("Next File to process: %v", *nextFile)

	err = s.ProcessTxnFile(nextFile)
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

	var nextFile os.DirEntry
	for _, file := range files {
		if !file.IsDir() {
			if s.fileNameRegExp.MatchString(file.Name()) {
				nextFile = file
				break
			}
		}
	}

	return &nextFile, nil
}

func (s DefaultService) ProcessTxnFile(txnFile *os.DirEntry) error {
	// Rename txn file while being processed
	originalFilename := (*txnFile).Name()
	processingFilename := "PROCESSING_" + originalFilename
	s.log.Infof("Original Name: %s", originalFilename)
	err := os.Rename(s.incomingDir+string(os.PathSeparator)+originalFilename, s.incomingDir+string(os.PathSeparator)+processingFilename)
	if err != nil {
		return err
	}

	// Open txn file
	f, err := os.Open(s.incomingDir + string(os.PathSeparator) + processingFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// do something with read line
		fmt.Printf("%+v\n", rec)
	}

	// Move txn file to archive
	err = os.Rename(s.incomingDir+string(os.PathSeparator)+processingFilename, s.archiveDir+string(os.PathSeparator)+originalFilename)
	if err != nil {
		return err
	}
	return nil
}
