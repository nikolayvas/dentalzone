package infrastructure

import (
	"crypto/tls"
	"strconv"

	config "dental_hub/configuration"

	gomail "gopkg.in/gomail.v2"
)

type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

func SendMail(mail *Mail) error {

	m := gomail.NewMessage()
	m.SetHeader("From", config.GetInstance().SMTP.Sender)
	m.SetHeader("To", mail.To[0])
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	host := config.GetInstance().SMTP.Host
	port, err := strconv.Atoi(config.GetInstance().SMTP.Port)
	if err != nil {
		return err
	}

	user := config.GetInstance().SMTP.Sender
	password := config.GetInstance().SMTP.Password

	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
