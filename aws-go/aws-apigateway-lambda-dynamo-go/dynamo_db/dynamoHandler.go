package dynamo_db

import (
	"github/dev-toolkit-go/aws-go/aws-apigateway-lambda-dynamo-go/model"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go/aws/session"
)

var handler DBHandler

type DynamoHandler struct {
	tableName   string
	dynamoDBAPI dynamodbiface.DynamoDBAPI
}

func NewDynamoHandler(tName string) DBHandler {
	if handler == nil {
		handler = &DynamoHandler{
			tableName:   tName,
			dynamoDBAPI: getDynamoInterface(),
		}
	}
	return handler
}

func getDynamoInterface() dynamodbiface.DynamoDBAPI {
	dynamoSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoInstance := dynamodb.New(dynamoSession)

	return dynamodbiface.DynamoDBAPI(dynamoInstance)
}

func convertToDBRecord(book *model.MyBook) map[string]*dynamodb.AttributeValue {
	item := map[string]*dynamodb.AttributeValue{
		"id":     {S: &book.ID},
		"title":  {S: &book.Title},
		"author": {S: &book.Author},
	}
	return item
}

func (h *DynamoHandler) Save(book *model.MyBook) error {
	input := &dynamodb.PutItemInput{
		Item:      convertToDBRecord(book),
		TableName: aws.String(h.tableName),
	}

	savedItem, err := h.dynamoDBAPI.PutItem(input)
	if err != nil {
		log.Fatal("Failed to save Item: ", err.Error())
		return err
	}

	log.Println("Item saved in db: ", savedItem)

	return nil
}

func (h *DynamoHandler) Update(book *model.MyBook) error {
	item, err := dynamodbattribute.MarshalMap(book)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(h.tableName),
	}

	updatedItem, err := h.dynamoDBAPI.PutItem(input)
	if err != nil {
		return err
	}

	log.Println("Item updated in db: ", updatedItem)
	return nil
}

func (h *DynamoHandler) UpdateAttributeByID(id, key, value string) error {
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				S: aws.String(value),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:        aws.String(h.tableName),
		UpdateExpression: aws.String("set " + key + " = :val"),
	}

	output, err := h.dynamoDBAPI.UpdateItem(&input)
	if err != nil {
		return err
	}

	log.Println("Item updated in db: ", output)

	return nil
}

func (h *DynamoHandler) GetByID(id string) (*model.MyBook, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(h.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}

	item, err := h.dynamoDBAPI.GetItem(input)
	if err != nil {
		return nil, err
	}

	if item.Item == nil {
		log.Fatal("Can't get item by id = ", id)
		return nil, nil
	}

	var book model.MyBook
	err = dynamodbattribute.UnmarshalMap(item.Item, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (h *DynamoHandler) DeleteByID(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(h.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}

	_, err := h.dynamoDBAPI.DeleteItem(input)
	if err != nil {
		log.Fatal("Can't delete item by id = ", id)
		return err
	}

	return nil
}
