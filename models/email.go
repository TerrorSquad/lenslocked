package models

import "github.com/wneessen/go-mail"

const (
	DefaultSender = "support@lenslocked.com"
)

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
		return nil, err
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
