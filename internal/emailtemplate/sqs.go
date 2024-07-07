package emailtemplate

import (
	"errors"
	"fmt"
	"stori-challenge/internal/sqsclient"
)

type SqsStoreEmailHandler struct {
	sqscli *sqsclient.SQSClient
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
