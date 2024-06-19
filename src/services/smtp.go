package services

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(email, subject, url string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL_USER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", `<h1>Email Confirmation</h1>
                    <h2>Hello `+email+`</h2>
                    <p>Thank you for joining us. Please confirm your email by clicking on the following link</p>
                    <a href='`+url+`'> Click here</a>
					atau masuk ke link `+url)

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_EMAIL_USER"), os.Getenv("SMTP_EMAIL_PASS"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
