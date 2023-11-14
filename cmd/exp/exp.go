package main

import (
	"github.com/joho/godotenv"
	"github.com/terrorsquad/lenslocked/models"
	"log"
	"os"
	"strconv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}
	es, err := models.NewEmailService(models.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}
	err = es.ForgotPassword("g.ninkovic@angeltech.rs", "https://lenslocked.com/reset?id=1234")
}
