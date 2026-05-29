package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SendHandler struct {
	producer *KafkaProducer
}

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := Message{
		UUID:    uuid.New().String(),
		From:    "Alice",
		To:      "Bob",
		Message: "Hello from Go producer",
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.producer.Send(ctx, msg); err != nil {
		log.Printf("ERROR send: %v", err)
		http.Error(w, "failed to send message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Message sent!")
}
