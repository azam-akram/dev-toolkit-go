package main

import (
	"context"
	"dev-toolkit-go/http-server-docker-go/handler"
	httpserver "dev-toolkit-go/http-server-docker-go/server"
	"log/slog"
	"net/http"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	ctx := context.Background()

	h := handler.NewRequestHandler(logger)

	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(h.HandleHealth))
	mux.Handle("/", http.HandlerFunc(h.HandleRequest))

	slog.Info("Starting HTTP server...")
	server := httpserver.NewServer("8080", logger, mux)
	if err := server.Start(ctx); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}

	slog.Info("HTTP server is running on port 8080")

	// Wait for a termination signal (e.g., SIGINT, SIGTERM) to gracefully shut down the server
	<-ctx.Done()

	slog.Info("Stopping HTTP server...")
	if err := server.Stop(ctx); err != nil {
		slog.Error("Failed to stop HTTP server", "error", err)
		os.Exit(1)
	}
}
