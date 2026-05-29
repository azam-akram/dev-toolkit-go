package utils

import (
	"context"
	"dev-toolkit-go/aws-lambda-sns-events-go/model"

	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func getAWSRegion() string {
	region := os.Getenv("AWS_REGION")
	log.Println("AWS region: ", region)
	return region
}

func getSNSTopicARN() string {
	topicArn := os.Getenv("SNS_TOPIC_ARN")
	log.Println("SNS Topic ARN: ", topicArn)
	return topicArn
}

func PublishEvent(ctx context.Context, event *model.Event) (msgId string, err error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(getAWSRegion()))
	if err != nil {
		return "", err
	}
	snsClient := sns.NewFromConfig(cfg)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	snsInput := &sns.PublishInput{
		Message:  aws.String(string(eventBytes)),
		TopicArn: aws.String(getSNSTopicARN()),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"name": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.Name),
			},
		},
	}

	snsMsg, err := snsClient.Publish(ctx, snsInput)
	if err != nil {
		return "", err
	}

	log.Println("Published event: ", snsMsg)
	return aws.ToString(snsMsg.MessageId), nil
}
