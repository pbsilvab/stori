package account

import (
	"fmt"
	"time"
)

const (
	RepositoryTypeDynamo = "dynamo"
)

type AccountInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Transaction struct {
	Account string
	ID      string
	Date    time.Time
	Amount  float64
}

type AccountTxRepository interface {
	SaveTransaction(transaction Transaction) error
}

type AccountInfoRepository interface {
	SaveAccountInfo(ai AccountInfo) (*AccountInfo, error)
	FindAccountInfo(id string) (*AccountInfo, error)
	List() (*[]AccountInfo, error)
}

func NewAccountInfoRepository(rt string) (AccountInfoRepository, error) {
	switch rt {
	case RepositoryTypeDynamo:
		return NewAccountInfoRepositoryDynamoDB(), nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", rt)
	}
}

func NewAccountTxRepository(rt string) (AccountTxRepository, error) {
	switch rt {
	case RepositoryTypeDynamo:
		return NewAccountTxRepositoryDynamoDB(), nil
	default:
		return nil, fmt.Errorf("unknown repository type: %s", rt)
	}
}
