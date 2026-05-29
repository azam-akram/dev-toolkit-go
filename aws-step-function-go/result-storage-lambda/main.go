package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type ProcessedData struct {
	ProcessedValue string `json:"processed_value"`
}

func HandleRequest(ctx context.Context, data ProcessedData) (string, error) {
	// Store data in database
	fmt.Printf("Storing result: %s\n", data.ProcessedValue)

	return "Data successfully stored!", nil
}

func main() {
	lambda.Start(HandleRequest)
}
