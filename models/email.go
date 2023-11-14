package models

import (
	"fmt"
	"github.com/wneessen/go-mail"
	html_template "html/template"
	text_template "text/template"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	To        string
	From      string
	Subject   string
	Plaintext string
	HTML      string
}
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) (*EmailService, error) {
	client, err := mail.NewClient(
		config.Host,
		mail.WithPort(config.Port),
		mail.WithSMTPAuth(mail.SMTPAuthCramMD5),
		mail.WithUsername(config.Username),
		mail.WithPassword(config.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating email client: %w", err)
	}
	es := EmailService{
		client: client,
	}
	return &es, nil
}

type EmailService struct {
	// DefaultSender is used as the default sender when one isn't provided for an
	// email. This is also used in functions where the email is predetermined,
	// like the forgotten password email.
	DefaultSender string
	// unexported fields
	client *mail.Client
}

func (emailService *EmailService) Send(email Email) error {
	msg := mail.NewMsg()
	emailService.setFrom(msg, email)
	msg.SetAddrHeader("To", email.To)
	msg.SetGenHeader("Subject", email.Subject)

	plainTextTemplate, _ := text_template.New("email").Parse(email.Plaintext)
	htmlTemplate, _ := html_template.New("email").Parse(email.HTML)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBodyTextTemplate(plainTextTemplate, nil)
		msg.AddAlternativeHTMLTemplate(htmlTemplate, nil)
	case email.Plaintext != "":
		msg.SetBodyTextTemplate(plainTextTemplate, nil)
	case email.HTML != "":
		msg.SetBodyHTMLTemplate(htmlTemplate, nil)
	}

	err := emailService.client.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}

func (emailService *EmailService) ForgotPassword(to string, resetURL string) error {
	email := Email{
		To:      to,
		Subject: "Reset your password",
		Plaintext: fmt.Sprintf(`To reset your password, use this link:
%s`, resetURL),
		HTML: fmt.Sprintf(`Click the link below to reset your password:<br>
<a href="%s">%s</a>`, resetURL, resetURL),
	}
	err := emailService.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
}

func (emailService *EmailService) setFrom(msg *mail.Msg, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case emailService.DefaultSender != "":
		from = emailService.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetAddrHeader("From", from)
}
