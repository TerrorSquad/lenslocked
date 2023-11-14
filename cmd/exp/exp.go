package main

import "github.com/terrorsquad/lenslocked/models"

const (
	port     = 2525
	host     = "hhhhhhhhhh"
	username = "uuuuuuuuuuu"
	password = "ppppppppppp"
)

func main() {
	es, err := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	email := models.Email{
		To:        "warhawk@hotmail.rs",
		From:      "g.ninkovic@angeltech.rs",
		Subject:   "Test email",
		Plaintext: "This is a test email",
		HTML:      "<h1>This is a test email</h1>",
	}
	if err := es.Send(email); err != nil {
		panic(err)
	}
}
