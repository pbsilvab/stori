package emailservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SES struct {
	senderEmail string
	sess        *ses.SES
}

func newSES(se string) SES {
	// Create a new AWS session using shared credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an SES service client
	svc := ses.New(sess)

	return SES{
		senderEmail: se,
		sess:        svc,
	}
}

func (svc SES) Send(content string, email string, subject string) (bool, error) {

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(email)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(content),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(svc.senderEmail),
	}

	_, err := svc.sess.SendEmail(input)

	if err != nil {
		return false, fmt.Errorf("error sending email: %v", err)
	}

	return true, nil
}
