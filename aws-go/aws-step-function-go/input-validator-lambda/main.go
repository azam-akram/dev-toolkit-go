package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Input struct {
	Value string `json:"value"`
}

type ValidationResponse struct {
	IsValid bool   `json:"is_valid"`
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, input Input) (ValidationResponse, error) {
	if input.Value == "" {
		return ValidationResponse{
			IsValid: false,
			Message: "Input value is empty.",
		}, nil
	}
	return ValidationResponse{
		IsValid: true,
		Message: "Input is valid.",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
