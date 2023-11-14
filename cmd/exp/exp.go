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
	err = es.ForgotPassword("g.ninkovic@angeltech.rs", "https://lenslocked.com/reset?id=1234")
}
