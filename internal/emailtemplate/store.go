package emailtemplate

import (
	"fmt"
	"os"
)

const (
	FsStorage = "fs"
	SQS       = "sqs"
)

type EmailTemplateStorage interface {
	StoreEmail(completedEmail string, email string) error
}

func NewStoreHandler(smt string) (EmailTemplateStorage, error) {
	switch smt {
	case FsStorage:
		outputDir := os.Getenv("EMAIL_FS_OUTPUT_DIR")
		if outputDir == "" {
			return nil, fmt.Errorf("for FS email storage you need to configure an output dir env var: %v", "EMAIL_FS_OUTPUT_DIR")
		}
		return NewFsEmailHandler(outputDir), nil
	case SQS:
		region := os.Getenv("SQS_REGION")
		if region == "" {
			return nil, fmt.Errorf("for SQS email storage you need to configure a sqs region env var: %v", "SQS_REGION")
		}

		sqsUrl := os.Getenv("SQS_URL")
		if sqsUrl == "" {
			return nil, fmt.Errorf("for SQS email storage you need to configure a sqs URL env var: %v", "SQS_URL")
		}

		sqsHandler, err := NewSQSEmailHandler(region, sqsUrl)
		if err != nil {
			return nil, fmt.Errorf("error loading sqs client: %v", err.Error())
		}

		return sqsHandler, nil
	default:
		return nil, fmt.Errorf("unknown email store handler type: %s", smt)
	}
}
