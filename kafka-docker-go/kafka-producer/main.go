package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "demo-kafka-topic")
	port := getEnv("PORT", "4444")

	producer := NewKafkaProducer(broker, topic)
	defer producer.Close()

	mux := http.NewServeMux()
	mux.Handle("/send", &SendHandler{producer: producer})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"up"}`))
	})

	log.Printf("Producer listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
