package main

import (
	"awesomeProject/config"
	"awesomeProject/model"
	"awesomeProject/service/mail"
)

func main() {
	// Logger
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	// Configs
	configs := config.LoadConfig(logger)

	// Services
	mailSvc := mail.NewDefaultService(logger, &configs.Smtp)

	var arrTxnCountPerMonth = []model.TxnCountPerMonth{
		{
			Month:    "January",
			TxnCount: 123,
		},
		{
			Month:    "February",
			TxnCount: 456,
		},
		{
			Month:    "March",
			TxnCount: 789,
		},
	}

	err := mailSvc.SendMail("gunterhm@gmail.com", arrTxnCountPerMonth)
	if err != nil {
		return
	}
}

/*func send(body string) {
	from := "gunterhm@gmail.com"
	pass := "rpnsgjmpgemiopna"
	to := "gunterhm@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
}*/
