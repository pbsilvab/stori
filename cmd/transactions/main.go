package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"stori-challenge/internal/account"
	"stori-challenge/internal/emailtemplate"
	"stori-challenge/internal/fileprocessor"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/joho/godotenv"
)

const (
	LAMBDA = "lambda"
	HTTP   = "http"
	CLI    = "cli"
)

var acctxr account.AccountTxRepository
var accinf account.AccountInfoRepository
var eth *emailtemplate.EmailTemplateHandler

type ProcessTxsLmbda struct {
	Directory string `json:"directory"`
	Account   string `json:"account"`
}

type TransactionParamsRequest struct {
	AccountID string `json:"accountId"`
	Directory string `json:"directory"`
}

func main() {
	godotenv.Load()
	runtime := os.Getenv("RUNTIME")
	var err error

	//Package Definitions

	repoType := os.Getenv("REPOSITORY_TYPE")
	// Account Transacctions Repository Definition
	acctxr, err = account.NewAccountTxRepository(repoType)
	if err != nil {
		log.Fatalf("error getting latest CSV file: %v", err)
	}

	// Account Information Repository Definition
	accinf, err = account.NewAccountInfoRepository(repoType)
	if err != nil {
		log.Fatalf("error creating account repository: %v", err)
	}

	// Email Storage Definition
	emailStorageHandler := os.Getenv("EMAIL_STORAGE_HANDLER_TYPE")
	esh, err := emailtemplate.NewStoreHandler(emailStorageHandler)
	if err != nil {
		log.Fatalf("error generating and saving email: %v", err.Error())
	}
	// Email Template Builder
	eth = emailtemplate.NewEmailTemplateHandler(esh)

	// Configure Runtime
	switch runtime {
	case HTTP:
		serveHttpApplication()
	case CLI:
		cliHandler()
	case LAMBDA:
		lambda.Start(ProcessTxsLambda)
	default:
		log.Fatal("No runtime defined")
	}
}

func ProcessTxsLambda(ctx context.Context, event ProcessTxsLmbda) (string, error) {
	accountID := event.Account
	directory := event.Directory

	err := processAccountTx(accountID, directory)

	if err != nil {
		return "", fmt.Errorf("error: %v", err.Error())
	}

	return "done", nil
}

func cliHandler() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <account_id> <directory>")
		os.Exit(1)
	}
	accountID := os.Args[1]
	directory := os.Args[2]

	err := processAccountTx(accountID, directory)

	if err != nil {
		log.Fatal(err)
	}
}

func serveHttpApplication() {
	http.HandleFunc("/process", httpHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST")
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var p TransactionParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if p.Directory == "" || p.AccountID == "" {
		http.Error(w, "directory, name and account_id are required", http.StatusBadRequest)
		return
	}

	err := processAccountTx(p.AccountID, p.Directory)

	if err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "ok",
	})
}

func processAccountTx(accountID string, directory string) error {

	fp := fileprocessor.NewFileProcessor(directory)

	records, err := fp.GetLatestCSVFile()
	if err != nil {
		return fmt.Errorf("error: %v ", err)
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

	accountInfo, err := accinf.FindAccountInfo(accountID)

	if err != nil {
		return fmt.Errorf("error find account info: %v", err)
	}

	acc := account.NewAccountTransactions(accountID, transactions, acctxr)
	totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth := acc.CalculateSummary()

	fmt.Printf("Account ID: %s\n", acc.ID)
	fmt.Printf("Balance: %.2f\n", totalBalance)

	//TODO: if saving is ok! Defer delete file, or store backup to s3
	acc.SaveTransactions()

	summaryContent := eth.GenerateSummaryContent(totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth)
	template := eth.GetDefaultTemplate()
	params := map[string]string{
		"Email":        accountInfo.Email,
		"Name":         accountInfo.Name,
		"TotalBalance": summaryContent,
	}

	if err := eth.GenerateAndSaveEmail(template, params, "tmp/emails"); err != nil {
		return fmt.Errorf("error generating and saving email: %v", err)
	}

	fmt.Println("Email generated and saved successfully.")

	return nil
}
