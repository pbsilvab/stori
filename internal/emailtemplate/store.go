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
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown email store handler type: %s", smt)
	}
}
