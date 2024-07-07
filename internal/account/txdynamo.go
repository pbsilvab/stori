package account

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoAccountTxRepository struct {
	tableName string
	svc       *dynamodb.DynamoDB
}

func NewAccountTxRepositoryDynamoDB() *DynamoAccountTxRepository {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return &DynamoAccountTxRepository{
		tableName: "transactions",
		svc:       svc,
	}
}

func (repo *DynamoAccountTxRepository) SaveTransaction(transaction Transaction) error {

	tx, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		"account": transaction.Account,
		"id":      transaction.ID,
		"date":    transaction.Date,
		"amount":  transaction.Amount,
	})

	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      tx,
		TableName: aws.String(repo.tableName),
	}

	_, err = repo.svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + transaction.ID + " to table " + repo.tableName)

	return nil
}
