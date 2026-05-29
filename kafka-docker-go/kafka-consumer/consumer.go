package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxRetries = 3
	retryDelay = time.Second
)

type KafkaConsumer struct {
	reader    *kafka.Reader
	dlqWriter *kafka.Writer
}

func NewKafkaConsumer(broker, topic, groupID string) *KafkaConsumer {
	dlqTopic := topic + "-dlt"
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    1,
			MaxBytes:    10e6,
			MaxWait:     time.Second,
			StartOffset: kafka.FirstOffset,
		}),
		dlqWriter: &kafka.Writer{
			Addr:                   kafka.TCP(broker),
			Topic:                  dlqTopic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

// Run fetches messages in a loop until ctx is cancelled.
func (c *KafkaConsumer) Run(ctx context.Context) {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // clean shutdown on SIGTERM / SIGINT
			}
			log.Printf("ERROR fetch: %v", err)
			continue
		}

		if err := c.processWithRetry(m); err != nil {
			log.Printf("ERROR retries exhausted, routing to DLQ — uuid=%s", string(m.Key))
			c.sendToDLQ(ctx, m)
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("ERROR commit: %v", err)
		}
	}
}

func (c *KafkaConsumer) processWithRetry(m kafka.Message) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := process(m); err != nil {
			lastErr = err
			log.Printf("WARN attempt %d/%d failed: %v", attempt, maxRetries, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}
			continue
		}
		return nil
	}
	return lastErr
}

func process(m kafka.Message) error {
	var msg Message
	if err := json.Unmarshal(m.Value, &msg); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	log.Println("=================================")
	log.Println("Received message:")
	log.Printf("  UUID:    %s", msg.UUID)
	log.Printf("  From:    %s", msg.From)
	log.Printf("  To:      %s", msg.To)
	log.Printf("  Message: %s", msg.Message)
	log.Println("=================================")
	return nil
}

func (c *KafkaConsumer) sendToDLQ(ctx context.Context, m kafka.Message) {
	if err := c.dlqWriter.WriteMessages(ctx, kafka.Message{
		Key:   m.Key,
		Value: m.Value,
	}); err != nil {
		log.Printf("ERROR DLQ write: %v", err)
	}
}

func (c *KafkaConsumer) Close() {
	c.reader.Close()
	c.dlqWriter.Close()
}
