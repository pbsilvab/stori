package sqsclient

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSClient represents the SQS client
type SQSClient struct {
	svc      *sqs.SQS
	queueURL string
}

// EmailMessage represents the structure of the email message
type EmailMessage struct {
	Template string `json:"template"`
	Email    string `json:"email"`
}

// New creates a new SQS client
func New(region, queueURL string) (*SQSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	svc := sqs.New(sess)
	return &SQSClient{svc: svc, queueURL: queueURL}, nil
}

// PushEmailMessage sends an email message to the SQS queue
func (c *SQSClient) PushEmailMessage(emailTemplate, email string) error {
	emailMessage := EmailMessage{
		Template: emailTemplate,
		Email:    email,
	}

	messageBody, err := json.Marshal(emailMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	_, err = c.svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// PollMessages receives messages from the SQS queue
func (c *SQSClient) PollMessages(maxMessages int64, waitTimeSeconds int64, visibilityTimeout int64) ([]*sqs.Message, error) {
	result, err := c.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: aws.Int64(maxMessages),
		WaitTimeSeconds:     aws.Int64(waitTimeSeconds),
		VisibilityTimeout:   aws.Int64(visibilityTimeout),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}
	return result.Messages, nil
}

// DeleteMessage deletes a message from the SQS queue
func (c *SQSClient) DeleteMessage(receiptHandle *string) error {
	_, err := c.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
