package emailtemplate

import (
	"fmt"
	"os"
	"path/filepath"
)

type FsEmailHandler struct {
	outputDir string
}

func NewFsEmailHandler(outputDir string) FsEmailHandler {
	return FsEmailHandler{
		outputDir: outputDir,
	}
}

func (fsh FsEmailHandler) StoreEmail(completedEmail string, email string) error {

	emailFileName := fmt.Sprintf("email_%v.txt", email)
	emailFilePath := filepath.Join(fsh.outputDir, emailFileName)

	if err := os.MkdirAll(fsh.outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	if err := os.WriteFile(emailFilePath, []byte(completedEmail), 0644); err != nil {
		return fmt.Errorf("error writing email file: %v", err)
	}

	return nil
}
