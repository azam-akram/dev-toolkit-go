package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Source string `json:"source,omitempty"`
	Action string `json:"action,omitempty"`
}

func HandleRequest(ctx context.Context, event Event) {
	log.Println("Context: ", ctx)
	log.Println("Event received: ", event)
}

func main() {
	lambda.Start(HandleRequest)
}
