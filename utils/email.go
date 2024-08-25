package utils

import (
	"os"
	"strconv"

	"github.com/vantutran2k1/social-network-auth/errors"
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) *errors.ApiError {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM_EMAIL"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), smtpPort, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
