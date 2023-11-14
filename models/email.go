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

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMsg()
	es.setFrom(msg, email)
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

	err := es.client.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}

func (es *EmailService) setFrom(msg *mail.Msg, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetAddrHeader("From", from)
}
