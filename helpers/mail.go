package helpers

import (
	gomail "gopkg.in/gomail.v2"
)

type MailConnection struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MailContent struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (conn *MailConnection) MailSend(
	from string,
	to string,
	subject string,
	body string,
) (bool, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(conn.Host, conn.Port, conn.Username, conn.Password)

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}
