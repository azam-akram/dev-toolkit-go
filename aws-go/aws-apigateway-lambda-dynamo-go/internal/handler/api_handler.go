package handler

import (
	"context"
	"encoding/json"
	"github/dev-toolkit-go/aws-go/aws-apigateway-lambda-dynamo-go/internal/dynamo_db"
	"github/dev-toolkit-go/aws-go/aws-apigateway-lambda-dynamo-go/internal/model"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
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
