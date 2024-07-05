package main

import (
	"context"
	"fmt"
	"log"
	"stori-challenge/internal/account"
	"stori-challenge/internal/emailtemplate"
	"stori-challenge/internal/fileprocessor"

	"github.com/aws/aws-lambda-go/lambda"
)

type ProcessTxs struct {
	Directory string `json:"directory"`
	Account   string `json:"account"`
	Name      string `json:"name"`
}

func ProcessTxsRequest(ctx context.Context, event ProcessTxs) (string, error) {
	accountID := event.Account
	name := event.Name
	directory := event.Directory

	fp := fileprocessor.NewFileProcessor(directory)

	records, err := fp.GetLatestCSVFile()
	if err != nil {
		log.Fatalf("Error getting latest CSV file: %v", err)
	}

	var transactions []account.Transaction
	for _, record := range records[1:] { // Skip header row
		transaction, err := account.ParseTransaction(record)
		if err != nil {
			log.Printf("Error parsing transaction: %v", err)
			continue
		}
		transactions = append(transactions, transaction)
	}

	acc := account.NewAccount(accountID, transactions)
	totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth := acc.CalculateSummary()

	fmt.Printf("Account ID: %s\n", acc.ID)
	fmt.Printf("Balance: %.2f\n", totalBalance)

	eth := emailtemplate.NewEmailTemplateHandler()

	summaryContent := eth.GenerateSummaryContent(totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth)

	template := eth.GetDefaultTemplate()

	params := map[string]string{
		"Name":         name,
		"TotalBalance": summaryContent,
	}

	if err := eth.GenerateAndSaveEmail(template, params, "tmp/emails"); err != nil {
		log.Fatalf("Error generating and saving email: %v", err)
	}

	return "done", nil
}

func main() {
	lambda.Start(ProcessTxsRequest)
}
