package model

// TxnCountPerMonth is the model for storing number of transactions per month for a given account
type TxnCountPerMonth struct {
	Month    string
	TxnCount int
}
