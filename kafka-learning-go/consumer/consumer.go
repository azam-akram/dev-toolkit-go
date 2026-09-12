package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	UUID    string `json:"uuid"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

const (
	maxRetries = 3
	retryDelay = time.Second
)

type Consumer struct {
	reader    *kafka.Reader
	dlqwriter *kafka.Writer
}

func NewConsumer(broker, topic, groupId string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			Topic:       topic,
			GroupID:     groupId,
			MinBytes:    10e3,
			MaxBytes:    10e6,
			StartOffset: kafka.FirstOffset,
		}),
		dlqwriter: &kafka.Writer{
			Addr:                   kafka.TCP(broker),
			Topic:                  topic + "-dlt",
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("ERROR fetch: %v", err)
			continue
		}

		if err := c.processWithRetry(m); err != nil {
			log.Printf("ERROR retries exhausted, routing to DLQ - uuid=%s", string(m.Key))
			c.sendToDLQ(ctx, m)
		}
	}
}

func (c *Consumer) processWithRetry(m kafka.Message) error {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := c.process(m); err != nil {
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

func (c *Consumer) process(m kafka.Message) error {
	var msg Message
	err := json.Unmarshal(m.Value, &msg)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
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

func (c *Consumer) sendToDLQ(ctx context.Context, m kafka.Message) {
	err := c.dlqwriter.WriteMessages(ctx, kafka.Message{
		Key:   m.Key,
		Value: m.Value,
	})

	if err != nil {
		log.Printf("ERROR DLQ write: %v", err)
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
	c.dlqwriter.Close()
}
