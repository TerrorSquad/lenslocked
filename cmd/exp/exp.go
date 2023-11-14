package main

import (
	"github.com/wneessen/go-mail"
	"log"
)

const (
	port     = 2525
	host     = "hhhhhhhhhh"
	username = "uuuuuuuuuuu"
	password = "ppppppppppp"
)

func main() {
	m := mail.NewMsg()
	if err := m.From("toni.sender@example.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To("tina.recipient@example.com"); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	m.Subject("This is my first mail with go-mail!")
	m.SetBodyString(mail.TypeTextPlain, "Do you like this mail? I certainly do!")
	c, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthCramMD5),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.Send(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
}
