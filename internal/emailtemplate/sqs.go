package emailtemplate

import (
	"encoding/json"
	"errors"
	"fmt"
	"stori-challenge/internal/sqsclient"
)

type SqsStoreEmailHandler struct {
	sqscli *sqsclient.SQSClient
}

type EmailMessage struct {
	Template string `json:"template"`
	Email    string `json:"email"`
}

func NewSQSEmailHandler(region string, sqsUrl string) (*SqsStoreEmailHandler, error) {
	sqscli, err := sqsclient.New(region, sqsUrl)

	if err != nil {
		msg := fmt.Sprintf("Got error SQS Email Handler: %s", err.Error())
		return nil, errors.New(msg)
	}

	return &SqsStoreEmailHandler{
		sqscli: sqscli,
	}, nil
}

func (svc SqsStoreEmailHandler) StoreEmail(completedEmail string, email string) error {

	emailMessage := EmailMessage{
		Template: completedEmail,
		Email:    email,
	}

	messageBody, err := json.Marshal(emailMessage)

	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}
	err = svc.sqscli.PushMessage(messageBody)

	return err
}
