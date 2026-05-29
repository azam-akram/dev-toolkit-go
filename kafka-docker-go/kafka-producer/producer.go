package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(broker),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			RequiredAcks:           kafka.RequireOne,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *KafkaProducer) Send(ctx context.Context, msg Message) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.UUID),
		Value: value,
	}); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Sent  uuid=%s topic=%s", msg.UUID, p.writer.Topic)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
