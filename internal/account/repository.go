package account

import (
	"fmt"
	"time"
)

const (
	RepositoryTypeDynamo = "dynamo"
)

type Transaction struct {
	Account string
	ID      string
	Date    time.Time
	Amount  float64
}

type AccountRepository interface {
	SaveTransaction(transaction Transaction) error
}

func NewAccountRepository(rt string) (AccountRepository, error) {
	switch rt {
	case RepositoryTypeDynamo:
		return NewAccountTxRepositoryDynamoDB(), nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", rt)
	}
}
