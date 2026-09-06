package helper

import (
	"fmt"
	"log"
	"server/config"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(cfg *config.Config, email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", email)
	m.SetHeader("From", cfg.EnvEmailUser)
	m.SetHeader("Subject", "Verify Your Music Room Account")

	url := fmt.Sprintf(
		"%s:%s/verify?token=%s",
		cfg.EnvAppHostname,
		cfg.EnvFrontendPort,
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
		cfg.EnvEmailUser,
		cfg.EnvEmailPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error when sending email:", err)
		return err
	}

	return nil
}

func SendPasswordResetEmail(cfg *config.Config, email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", email)
	m.SetHeader("From", cfg.EnvEmailUser)
	m.SetHeader("Subject", "Reset Your Music Room Password")

	url := fmt.Sprintf(
		"%s:%s/resetpassword?token=%s",
		cfg.EnvAppHostname,
		cfg.EnvFrontendPort,
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
		cfg.EnvEmailUser,
		cfg.EnvEmailPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error when sending email:", err)
		return err
	}

	return nil
}
