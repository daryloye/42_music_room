package helper

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", email)
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("Subject", "Verify Your Music Room Account")

	url := fmt.Sprintf(
		"%s:%s/verify?token=%s",
		os.Getenv("APP_HOSTNAME"),
		os.Getenv("FRONTEND_PORT"),
		token,
	)

	message := fmt.Sprintf(
		`<h2>Hello!</h2>
		<p>Please verify your Music Room account by clicking this link:</p>
		<a href="%s">Verify account</a>`,
		url,
	)

	m.SetBody("text/html", message)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error when sending email:", err)
		return err
	}

	return nil
}

func SendPasswordResetEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", email)
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("Subject", "Reset Your Music Room Password")

	url := fmt.Sprintf(
		"%s:%s/resetpassword?token=%s",
		os.Getenv("APP_HOSTNAME"),
		os.Getenv("FRONTEND_PORT"),
		token,
	)

	message := fmt.Sprintf(
		`<h2>Hello!</h2>
		<p>Please reset your Music Room password by clicking this link:</p>
		<a href="%s">Reset password</a>
		<p>The link will expire in 15 minutes</p>`,
		url,
	)

	m.SetBody("text/html", message)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error when sending email:", err)
		return err
	}

	return nil
}
