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

type ProcessTxsLmbda struct {
	Directory string `json:"directory"`
	Account   string `json:"account"`
	Name      string `json:"name"`
}

type TransactionParamsRequest struct {
	AccountID string `json:"accountId"`
	Name      string `json:"name"`
	Directory string `json:"directory"`
}

func main() {
	godotenv.Load()
	runtime := os.Getenv("RUNTIME")

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
	name := event.Name
	directory := event.Directory

	err := processAccountTx(accountID, name, directory)

	if err != nil {
		return "", fmt.Errorf("error: %v", err.Error())
	}

	return "done", nil
}

func cliHandler() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: main <account_id> <name> <directory>")
		os.Exit(1)
	}
	accountID := os.Args[1]
	name := os.Args[2]
	directory := os.Args[3]

	err := processAccountTx(accountID, name, directory)

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

	if p.Directory == "" || p.AccountID == "" || p.Name == "" {
		http.Error(w, "directory, name and account_id are required", http.StatusBadRequest)
		return
	}

	err := processAccountTx(p.AccountID, p.Name, p.Directory)

	if err != nil {
		http.Error(w, "error decoding request body: "+err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "ok",
	})
}

func processAccountTx(accountID string, name string, directory string) error {

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

	repoType := os.Getenv("REPOSITORY_TYPE")
	r, err := account.NewAccountRepository(repoType)

	if err != nil {
		return fmt.Errorf("error getting latest CSV file: %v", err)
	}

	acc := account.NewAccount(accountID, transactions, r)
	totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth := acc.CalculateSummary()

	fmt.Printf("Account ID: %s\n", acc.ID)
	fmt.Printf("Balance: %.2f\n", totalBalance)

	//TODO: if saving is ok! Defer delete file, or store backup to s3
	acc.SaveTransactions()

	eth := emailtemplate.NewEmailTemplateHandler()
	summaryContent := eth.GenerateSummaryContent(totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth)
	template := eth.GetDefaultTemplate()
	params := map[string]string{
		"Name":         name,
		"TotalBalance": summaryContent,
	}

	if err := eth.GenerateAndSaveEmail(template, params, "tmp/emails"); err != nil {
		return fmt.Errorf("error generating and saving email: %v", err)
	}

	fmt.Println("Email generated and saved successfully.")

	return nil
}
