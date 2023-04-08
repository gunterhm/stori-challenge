package mail

import (
	"awesomeProject/config"
	"awesomeProject/model"
	"bufio"
	"bytes"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"html/template"
)

// DefaultService is an implementation of IMailService
type DefaultService struct {
	log        *zap.SugaredLogger
	smtpConfig *config.SmptConfigurations
}

// NewDefaultService creates a new DefaultService
func NewDefaultService(logger *zap.SugaredLogger, smtpConf *config.SmptConfigurations) *DefaultService {
	return &DefaultService{
		log:        logger,
		smtpConfig: smtpConf,
	}
}

func (s DefaultService) SendSummaryMail(summaryEmailData *model.SummaryEmail) error {
	tmpl, err := template.ParseFiles("resources/email_templates/summary_email.html")
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	err = tmpl.Execute(writer, *summaryEmailData)
	if err != nil {
		return err
	}

	err = writer.Flush() // forcefully write remaining
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", s.smtpConfig.From)
	mail.SetHeader("To", summaryEmailData.Email)
	mail.SetHeader("Subject", "Account Summary")
	mail.SetBody("text/html", buffer.String())
	mail.Embed("resources/stori-logo.png")

	s.log.Infof("EMail: %s", buffer.String())

	if s.smtpConfig.SendMail {
		d := gomail.NewDialer(s.smtpConfig.Host, s.smtpConfig.Port, s.smtpConfig.User, s.smtpConfig.Password)

		// Send the email
		if err := d.DialAndSend(mail); err != nil {
			return err
		}
	}

	return nil
}
