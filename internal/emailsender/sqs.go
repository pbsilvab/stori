package emailsender

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"stori-challenge/internal/emailservice"
	"stori-challenge/internal/sqsclient"
)

type SqsStoreEmailList struct {
	sqscli *sqsclient.SQSClient
	sender emailservice.EmailService
}

func NewSQSEmailList(region string, sqsUrl string, sender emailservice.EmailService) (*SqsStoreEmailList, error) {
	sqscli, err := sqsclient.New(region, sqsUrl)

	if err != nil {
		msg := fmt.Sprintf("Got error SQS Email Handler: %s", err.Error())
		return nil, errors.New(msg)
	}

	return &SqsStoreEmailList{
		sqscli: sqscli,
		sender: sender,
	}, nil
}

func (svc SqsStoreEmailList) Pull() (*EmailContent, error) {
	// This could be better :)

	messges, err := svc.sqscli.PollMessages(1, 10, 30)
	if err != nil {
		return nil, fmt.Errorf("error parsing message: %v", err)
	}

	if len(messges) < 1 {
		return nil, fmt.Errorf("no message in the queue: %v", "Error")
	}

	var parsedMessages []EmailContent
	for _, sqsMsg := range messges {
		var msg EmailContent
		err := json.Unmarshal([]byte(*sqsMsg.Body), &msg)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}
		parsedMessages = append(parsedMessages, msg)

		svc.sqscli.DeleteMessage(sqsMsg.ReceiptHandle)
	}

	return &parsedMessages[0], nil
}

func (svc SqsStoreEmailList) Send(content *EmailContent) error {
	svc.sender.Send(content.Template, content.Email, "Monthly Account Summary SVG")
	return nil
}

func (svc SqsStoreEmailList) Delete(params EmailContent) error {

	return nil
}
