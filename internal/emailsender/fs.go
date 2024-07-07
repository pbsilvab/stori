package emailsender

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"stori-challenge/internal/emailservice"
)

type FsEmailSender struct {
	inputDir string
	sender   emailservice.EmailService
}

func NewFsEmailSender(inputDir string, sender emailservice.EmailService) FsEmailSender {
	return FsEmailSender{
		inputDir: inputDir,
		sender:   sender,
	}
}

func (fs FsEmailSender) Pull() (*EmailContent, error) {
	files, err := os.ReadDir(fs.inputDir)

	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in directory")
	}

	lastFile := files[len(files)-1]
	fileName := lastFile.Name()
	email, err := extractEmailFromFileName(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to extract email from filename: %w", err)
	}

	filePath := filepath.Join(fs.inputDir, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &EmailContent{
		Template: string(content),
		Email:    email,
	}, nil
}

func (fs FsEmailSender) Send(content *EmailContent) error {
	fs.sender.Send(content.Template, content.Email, "default subject")
	return nil
}

func (fs FsEmailSender) Delete(content EmailContent) error {

	return nil
}

func extractEmailFromFileName(fileName string) (string, error) {
	regex := regexp.MustCompile(`email_(.+?)@gmail\.com\.txt`)
	matches := regex.FindStringSubmatch(fileName)
	if len(matches) < 2 {
		return "", fmt.Errorf("no email found in filename")
	}
	return matches[1] + "@gmail.com", nil
}
