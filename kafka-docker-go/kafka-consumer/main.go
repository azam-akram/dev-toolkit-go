package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	broker := "localhost:9092"
	topic := "demo-kafka-topic"
	groupID := "demo-kafka-consumer"
	port := "5555"

	consumer := NewKafkaConsumer(broker, topic, groupID)
	defer consumer.Close()

	// Simple health endpoint
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"up"}`))
		})
		log.Printf("Health endpoint on :%s", port)
		http.ListenAndServe(":"+port, mux)
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("Consumer started topic=%s group=%s", topic, groupID)
	consumer.Run(ctx)
	log.Println("Consumer stopped")
}
