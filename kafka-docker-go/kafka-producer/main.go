package main

import (
	"log"
	"net/http"
)

func main() {
	broker := "localhost:9092"
	topic := "demo-kafka-topic"
	port := "4444"

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
