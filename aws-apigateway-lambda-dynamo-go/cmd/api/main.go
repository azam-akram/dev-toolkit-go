package main

import (
	"dev-toolkit-go/aws-apigateway-lambda-dynamo-go/internal/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.HandleRequest)
}
