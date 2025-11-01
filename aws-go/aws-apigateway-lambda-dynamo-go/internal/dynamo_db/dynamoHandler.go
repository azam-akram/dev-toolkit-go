package dynamo_db

import (
	"context"
	"github/dev-toolkit-go/aws-go/aws-apigateway-lambda-dynamo-go/internal/model"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var handler DBHandler

type DynamoHandler struct {
	tableName      string
	dynamoDBClient *dynamodb.Client
}

func NewDynamoHandler(tName string) DBHandler {
	if handler == nil {
		handler = &DynamoHandler{
			tableName:      tName,
			dynamoDBClient: getDynamoClient(),
		}
	}
	return handler
}

func getDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func convertToDBRecord(book *model.MyBook) map[string]types.AttributeValue {
	item := map[string]types.AttributeValue{
		"id":     &types.AttributeValueMemberS{Value: book.ID},
		"title":  &types.AttributeValueMemberS{Value: book.Title},
		"author": &types.AttributeValueMemberS{Value: book.Author},
	}
	return item
}

func (h *DynamoHandler) Save(book *model.MyBook) error {
	input := &dynamodb.PutItemInput{
		Item:      convertToDBRecord(book),
		TableName: aws.String(h.tableName),
	}

	savedItem, err := h.dynamoDBClient.PutItem(context.TODO(), input)
	if err != nil {
		log.Fatal("Failed to save Item: ", err.Error())
		return err
	}

	log.Println("Item saved in db: ", savedItem)

	return nil
}

func (h *DynamoHandler) Update(book *model.MyBook) error {
	item, err := attributevalue.MarshalMap(book)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(h.tableName),
	}

	updatedItem, err := h.dynamoDBClient.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Println("Item updated in db: ", updatedItem)
	return nil
}

func (h *DynamoHandler) UpdateAttributeByID(id, key, value string) error {
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberS{Value: value},
		},
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		TableName:        aws.String(h.tableName),
		UpdateExpression: aws.String("set " + key + " = :val"),
	}

	output, err := h.dynamoDBClient.UpdateItem(context.TODO(), &input)
	if err != nil {
		return err
	}

	log.Println("Item updated in db: ", output)

	return nil
}

func (h *DynamoHandler) GetByID(id string) (*model.MyBook, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(h.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	item, err := h.dynamoDBClient.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if item.Item == nil {
		log.Println("Can't get item by id =", id)
		return nil, nil
	}

	var book model.MyBook
	err = attributevalue.UnmarshalMap(item.Item, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (h *DynamoHandler) DeleteByID(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(h.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := h.dynamoDBClient.DeleteItem(context.TODO(), input)
	if err != nil {
		log.Fatal("Can't delete item by id = ", id)
		return err
	}

	return nil
}
