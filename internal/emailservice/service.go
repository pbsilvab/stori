package emailservice

type EmailService interface {
	Send(content string, email string, subject string) (bool, error)
}

func NewEmailService(senderEmail string) EmailService {
	return newSES(senderEmail)
}
