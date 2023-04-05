package mail

import "awesomeProject/model"

// IMailService is a service interface for smtp services
type IMailService interface {
	SendMail(to string, txnCntPerMonth []model.TxnCountPerMonth) error
}
