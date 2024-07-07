package emailtemplate

import (
	"fmt"
	"strings"
)

const summaryTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Template</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #003A40;
            font-family: Arial, sans-serif;
            color: #ffffff;
        }
        .container {
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            color: #000000;
        }
        .banner {
            width: 100%;
            background-color: #003A40;
            text-align: center;
            padding: 20px 0;
        }
        .banner img {
            max-width: 100%;
            height: auto;
        }
        .content {
            padding: 20px;
        }
        .content p {
            font-size: 16px;
            line-height: 1.5;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="banner">
            <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTqgUzyqW16XTfgqiucKIak04aDMWel1rG8P0NFdhpylnjlKdt-K351AA1M6XS9EFqxg8k&usqp=CAU">
        </div>
        <div class="content">
            <p>
				Dear {{Name}},
			</p>
			<p>
				Thank you for your recent transaction with us. Here is the summary of your transactions:
			</p>
			<p>
				Total balance is {{TotalBalance}}
			</p>
				
				{{Transactions}}

				Regards,
				Stori
			</p>
        </div>
    </div>
</body>
</html>`

// EmailTemplateHandler represents an email template handler.
type EmailTemplateHandler struct {
	storeEmailSvc EmailTemplateStorage
	// You can add fields here to manage template data or configuration
}

// NewEmailTemplateHandler creates a new instance of EmailTemplateHandler.
func NewEmailTemplateHandler(svc EmailTemplateStorage) *EmailTemplateHandler {
	return &EmailTemplateHandler{
		storeEmailSvc: svc,
	}
}

// GetDefaultTemplate returns the default email template as a string.
func (eth *EmailTemplateHandler) GetDefaultTemplate() string {
	return summaryTemplate
	// 	return `Dear {{Name}},

	// Thank you for your recent transaction with us. Here is the summary of your transactions:

	// Total balance is {{TotalBalance}}
	// {{Transactions}}

	// Regards,
	// Stori`
}

func (eth *EmailTemplateHandler) GenerateSummaryContent(totalBalance float64, transactionsByMonth map[string]int, averageCreditByMonth, averageDebitByMonth map[string]float64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total balance is %.2f\n", totalBalance))
	for month, count := range transactionsByMonth {
		sb.WriteString(fmt.Sprintf("Number of transactions in %s: %d\n", month, count))
		sb.WriteString(fmt.Sprintf("Average credit amount in %s: %.2f\n", month, averageCreditByMonth[month]))
		sb.WriteString(fmt.Sprintf("Average debit amount in %s: %.2f\n", month, averageDebitByMonth[month]))
	}
	return sb.String()
}

func (eth *EmailTemplateHandler) GenerateAndSaveEmail(template string, params map[string]string, outputDir string) error {
	completedEmail := eth.populateTemplate(template, params)

	err := eth.storeEmailSvc.StoreEmail(completedEmail, params["Email"])

	if err != nil {
		return fmt.Errorf("error writing email file: %v", err.Error())
	}

	return nil
}

func (eth *EmailTemplateHandler) populateTemplate(template string, params map[string]string) string {
	completedEmail := template

	var sb strings.Builder
	for month, count := range params["transactionsByMonth"] {
		sb.WriteString(fmt.Sprintf("Number of transactions in %s: %d\n", fmt.Sprint(month), count))
	}
	params["Transactions"] = sb.String()

	for key, value := range params {
		placeholder := fmt.Sprintf("{{%s}}", key)
		completedEmail = strings.ReplaceAll(completedEmail, placeholder, value)
	}
	return completedEmail
}
