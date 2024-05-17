package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github/GoDevKit/awsGoDevKit/aws-lambda-external-sns-topic-go/calculation-requester-lambda/model"
	"github/GoDevKit/awsGoDevKit/aws-lambda-external-sns-topic-go/calculation-requester-lambda/utils"

	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) error {
	log.Println("SNS event received:", snsEvent)

	for _, record := range snsEvent.Records {
		var event model.Event
		err := json.Unmarshal([]byte(record.SNS.Message), &event)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			return err
		}

		switch event.Name {
		case "SumCompleted":
			if event.Source == "Calculator" {
				log.Println("Answer received:", event.Payload.Sum)
				return nil
			}

		case "StartingEvent":
			event.Name = "SumRequested"
			event.Source = "Calculation Requester"
			event.EventTime = time.Now().Format(time.RFC3339)

			log.Println("Event to publish: ", event)

			if _, err := utils.PublishEvent(context.Background(), event); err != nil {
				return fmt.Errorf("error publishing event: %w", err)
			}

			return nil

		default:
			log.Println("Unknown event, ignoring this..")
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
