package main

import (
	"fmt"
	"log"
	"os"
	"stori-challenge/internal/emailsender"
	"stori-challenge/internal/emailservice"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	senderEmail := os.Getenv("AWS_SES_SENDER_EMAIL")
	if senderEmail == "" {
		log.Fatalf("missing env: %v", "AWS_SES_SENDER_EMAIL")
	}

	esvc := emailservice.NewEmailService(senderEmail)

	storage := os.Getenv("EMAIL_STORAGE_HANDLER_TYPE")
	fmt.Println("storgeType", storage)
	if storage == "" {
		log.Fatalf("missing env: %v", "EMAIL_STORAGE_HANDLER_TYPE")
	}

	sender, err := emailsender.NewEmailSender(storage, esvc)

	if err != nil {
		log.Fatal(err)
	}

	ec, err := sender.Pull()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("start send")
	err = sender.Send(ec)
	fmt.Println("end send")

	if err != nil {
		log.Fatal(err)
	}
}
