package main

import (
	"context"
	"github/dev-toolkit-go/aws-go/apigateway-lambda-dynamo-aws-go/dynamo_db"
	"github/dev-toolkit-go/aws-go/apigateway-lambda-dynamo-aws-go/model"

	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, book model.MyBook) (*model.MyBook, error) {
	log.Println("Received event with Book: ", book)

	dynamoHandler := dynamo_db.NewDynamoHandler()
	if err := dynamoHandler.Save(&book); err != nil {
		log.Fatal("Failed to save item, error: ", err.Error())
	}

	if err := dynamoHandler.UpdateAttributeByID(book.ID, "author", "Modified Author"); err != nil {
		log.Fatal("Failed to update item's value by ID, error: ", err.Error())
	}

	updatedBook, err := dynamoHandler.GetByID(book.ID)
	if err != nil {
		log.Fatal("Failed to get item by ID, error: ", err.Error())
	}

	log.Println("Fetched updated Book from db: ", updatedBook)

	err = dynamoHandler.DeleteByID(book.ID)
	if err != nil {
		log.Fatal("Failed to delete item by ID, error: ", err.Error())
	}

	log.Println("Item deleted")

	return updatedBook, nil
}

func main() {
	lambda.Start(HandleRequest)
}
