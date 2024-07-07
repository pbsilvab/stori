package emailsender

import (
	"fmt"
	"os"

	"stori-challenge/internal/emailservice"
)

type EmailContent struct {
	Template string `json:"template"`
	Email    string `json:"email"`
}

const (
	FsStorage = "fs"
	SQS       = "sqs"
)

type EmailSender interface {
	Pull() (*EmailContent, error)
	Send(params *EmailContent) error
	Delete(params EmailContent) error
}

func NewEmailSender(handlerType string, emailService emailservice.EmailService) (EmailSender, error) {
	switch handlerType {
	case FsStorage:
		inputDir := os.Getenv("EMAIL_FS_INPUT_DIR")
		if inputDir == "" {
			return nil, fmt.Errorf("for FS email storage you need to configure an output dir env var: %v", "EMAIL_FS_INPUT_DIR")
		}
		return NewFsEmailSender(inputDir, emailService), nil
	case SQS:
		region := os.Getenv("SQS_REGION")
		if region == "" {
			return nil, fmt.Errorf("for SQS email storage you need to configure a sqs region env var: %v", "SQS_REGION")
		}

		sqsUrl := os.Getenv("SQS_URL")
		if sqsUrl == "" {
			return nil, fmt.Errorf("for SQS email storage you need to configure a sqs URL env var: %v", "SQS_URL")
		}

		sqsHandler, err := NewSQSEmailList(region, sqsUrl, emailService)
		if err != nil {
			return nil, fmt.Errorf("error loading sqs client: %v", err.Error())
		}

		return sqsHandler, nil
	default:
		return nil, fmt.Errorf("unknown email store handler type: %s", handlerType)
	}
}
