package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Input struct {
	Value string `json:"value"`
}

type ProcessedData struct {
	ProcessedValue string `json:"processed_value"`
}

func HandleRequest(ctx context.Context, input Input) (ProcessedData, error) {
	processedValue := fmt.Sprintf("Processed_%s", input.Value)

	return ProcessedData{
		ProcessedValue: processedValue,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
