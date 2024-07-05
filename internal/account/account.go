package account

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Account represents an account with transactions.
type Account struct {
	ID           string
	Transactions []Transaction
	repository   AccountRepository
}

// NewAccount creates a new account.
func NewAccount(id string, transactions []Transaction, repository AccountRepository) *Account {
	return &Account{
		ID:           id,
		Transactions: transactions,
		repository:   repository,
	}
}

func (a *Account) SaveTransactions() {
	for _, transaction := range a.Transactions {
		transaction.Account = a.ID
		a.repository.SaveTransaction(transaction)
	}
}

// CalculateSummary calculates the summary for the account.
func (a *Account) CalculateSummary() (float64, map[string]int, map[string]float64, map[string]float64) {
	totalBalance := 0.0
	transactionsByMonth := make(map[string]int)
	totalCreditByMonth := make(map[string]float64)
	totalDebitByMonth := make(map[string]float64)

	for _, transaction := range a.Transactions {
		totalBalance += transaction.Amount
		month := transaction.Date.Format("January")
		transactionsByMonth[month]++
		if transaction.Amount > 0 {
			totalCreditByMonth[month] += transaction.Amount
		} else {
			totalDebitByMonth[month] += transaction.Amount
		}
	}

	averageCreditByMonth := make(map[string]float64)
	averageDebitByMonth := make(map[string]float64)
	for month, count := range transactionsByMonth {
		if count > 0 {
			averageCreditByMonth[month] = totalCreditByMonth[month] / float64(count)
			averageDebitByMonth[month] = totalDebitByMonth[month] / float64(count)
		}
	}

	return totalBalance, transactionsByMonth, averageCreditByMonth, averageDebitByMonth
}

// ParseTransaction parses a transaction record and returns a Transaction.
func ParseTransaction(record []string) (Transaction, error) {
	id := record[0]
	date, err := time.Parse("1/2", record[1])
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid date format: %v", err)
	}
	amount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		if strings.HasPrefix(record[2], "+") {
			amount, _ = strconv.ParseFloat(record[2][1:], 64)
		} else {
			return Transaction{}, fmt.Errorf("invalid amount format: %v", err)
		}
	}

	return Transaction{
		ID:     id,
		Date:   date,
		Amount: amount,
	}, nil
}
