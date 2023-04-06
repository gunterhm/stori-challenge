package mail

import "awesomeProject/model"

// IMailService is a service interface for smtp services
type IMailService interface {
	SendSummaryMail(summaryEmailData *model.SummaryEmail) error
}
