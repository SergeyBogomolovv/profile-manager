package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/SergeyBogomolovv/profile-manager/notification/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendLoginEmail(to string, data domain.LoginNotification) error
	SendRegisterEmail(to string) error
}

type mailer struct {
	user string
	pass string
	port int
	host string
}

func New(cfg config.SMTP) Mailer {
	return &mailer{
		user: cfg.User,
		pass: cfg.Pass,
		port: cfg.Port,
		host: cfg.Host,
	}
}

func (m *mailer) SendLoginEmail(to string, data domain.LoginNotification) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.user)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Произведен вход в аккаунт")

	tmpl, err := template.ParseFiles("templates/login_notification.html")
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	mail.SetBody("text/html", body.String())

	d := gomail.NewDialer(m.host, m.port, m.user, m.pass)
	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("failed to send login email: %w", err)
	}
	return nil
}

func (m mailer) SendRegisterEmail(to string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.user)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Добро пожаловать в profile-manager")

	tmpl, err := template.ParseFiles("templates/register_notification.html")
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, nil); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	mail.SetBody("text/html", body.String())

	d := gomail.NewDialer(m.host, m.port, m.user, m.pass)
	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("failed to send register email: %w", err)
	}
	return nil
}
