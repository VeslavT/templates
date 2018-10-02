package golang

import (
	"net/smtp"
	"strings"
)

const (
	MailHeaders = "MailHeaders-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	MailHost    = "localhost:25"
)

func SendMail(from string, to []string, subject, body string) error {
	c, err := smtp.Dial(MailHost)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}

	mailBody := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" + MailHeaders + "\r\n" + body)

	_, err = w.Write(mailBody)

	if err != nil {
		return err
	}

	err = w.Close()

	if err != nil {
		return err
	}

	return c.Quit()
}

