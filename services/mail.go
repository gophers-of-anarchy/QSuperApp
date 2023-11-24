package services

import (
	"fmt"
	"net/smtp"
)

type EmailRequest struct {
	Username   string `json:"username"`
	SecretCode string `json:"secretCode"`
}

type EmailResponse struct {
	IsVerified bool   `json:"is-verified"`
	Message    string `json:"message"`
}

var (
	host     = "smtp.gmail.com"
	port     = "587"
	from     = "Q_APP@gmail.com"
	password = "zlafvgphwezlklai"
	//password = os.Getenv("GmailPassword")
)

func SendMail(subject, body string, to []string) error {
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(fmt.Sprint(host+":"+port), auth, from, to, []byte(fmt.Sprint("Subject: "+subject+"\n"+body)))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
