package main

import (
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/handler"
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/logger"

	"net/http"
)

type Student struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func main() {
	logLevel := "INFO"
	log := logger.Init(logLevel)

	http.HandleFunc("/", handler.HandleRequest)

	log.Info("HTTP server starting on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("Failed to start http server", "error", err)
	}

	log.Info("HTTP server is listening on :8080")
}
