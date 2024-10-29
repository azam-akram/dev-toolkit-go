package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dev-toolkit-go/aws-apigateway-lambda-demo-go/dynamo_db"
	"github.com/dev-toolkit-go/aws-apigateway-lambda-demo-go/model"
)

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoHandler := dynamo_db.NewDynamoHandler("my-demo-dynamo-table")

	switch req.HTTPMethod {
	case "POST":
		var book model.MyBook
		if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
		if err := dynamoHandler.Save(&book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: req.Body}, nil

	case "GET":
		id := req.QueryStringParameters["id"]
		book, err := dynamoHandler.GetByID(id)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		if book == nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
		}
		body, err := json.Marshal(book)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(body)}, nil

	case "PUT":
		var book model.MyBook
		if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
		if err := dynamoHandler.Update(&book); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: req.Body}, nil

	case "DELETE":
		id := req.QueryStringParameters["id"]
		if err := dynamoHandler.DeleteByID(id); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil

	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func main() {
	lambda.Start(HandleRequest)
}
