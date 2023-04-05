package model

// SummaryEmail is the model for storing data that will be used to fill in summary email template
type SummaryEmail struct {
	AccountNumber       string
	TotalBalance        float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	ArrTxnCountPerMonth []TxnCountPerMonth
}

// TxnCountPerMonth is the model for storing number of transactions per month for a given account
type TxnCountPerMonth struct {
	Month    string
	TxnCount int
}
