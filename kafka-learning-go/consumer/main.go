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
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "demo-kafka-topic")
	groupID := getEnv("KAFKA_GROUP_ID", "demo-kafka-consumer")
	port := getEnv("PORT", "5555")

	consumer := NewConsumer(broker, topic, groupID)
	defer consumer.Close()

	// healthcheck endpoint
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "OK"}`))
		})
		log.Printf("Running health endpoint on %s", port)
		http.ListenAndServe(":"+port, mux)
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("Consumer started topic=%s group=%s", topic, groupID)
	consumer.Run(ctx)
	log.Println("Consumer stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); len(v) > 1 {
		return v
	}
	return fallback
}
