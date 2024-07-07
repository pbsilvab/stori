package account

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Account represents an account with transactions.
type AccInfo struct {
	tableName string
	svc       *dynamodb.DynamoDB
}

func NewAccountInfoRepositoryDynamoDB() *AccInfo {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return &AccInfo{
		tableName: "accountInfo",
		svc:       svc,
	}
}

func (aci *AccInfo) SaveAccountInfo(ai AccountInfo) (*AccountInfo, error) {
	id, err := getID(12)
	if err != nil {
		msg := fmt.Sprintf("Got error gen id: %s", err.Error())
		return nil, errors.New(msg)
	}
	ai.Id = id
	item, err := dynamodbattribute.MarshalMap(ai)

	if err != nil {
		msg := fmt.Sprintf("Got error marshalling map: %s", err.Error())
		return nil, errors.New(msg)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(aci.tableName),
	}

	_, err = aci.svc.PutItem(input)
	if err != nil {
		msg := fmt.Sprintf("Got error calling PutItem: %s", err.Error())
		return nil, errors.New(msg)
	}

	fmt.Println("Successfully added '" + ai.Id + " to table " + aci.tableName)

	return &AccountInfo{
		Id:    ai.Id,
		Name:  ai.Name,
		Email: ai.Email,
	}, nil
}

func (aci *AccInfo) FindAccountInfo(id string) (*AccountInfo, error) {

	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(aci.tableName),
		Key:       key,
	}

	result, err := aci.svc.GetItem(input)

	if err != nil {
		msg := "Got error calling GetItem: " + err.Error()
		return nil, errors.New(msg)
	}

	if result.Item == nil {
		msg := "Could not find '" + id + "'"
		return nil, errors.New(msg)
	}

	item := AccountInfo{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &item, nil
}

func (aci *AccInfo) List() (*[]AccountInfo, error) {
	// Create the input for Scan
	input := &dynamodb.ScanInput{
		TableName: &aci.tableName,
	}

	// Perform the scan operation
	result, err := aci.svc.Scan(input)
	if err != nil {
		log.Fatalf("Got error scanning table: %s", err)
	}

	// Unmarshal the scan result into a slice of User structs
	var accs []AccountInfo
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &accs)
	if err != nil {
		log.Fatalf("Failed to unmarshal scan result: %s", err)
	}

	return &accs, nil
}

func getID(length int) (string, error) {
	bytes := make([]byte, length/2) // Since each byte will be represented by two hex characters
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
