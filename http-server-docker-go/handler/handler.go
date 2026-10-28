package handler

import (
	"log/slog"
	"net/http"
)

type Handler interface {
	HandleHealth(w http.ResponseWriter, r *http.Request)
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

type requestHandler struct {
	logger *slog.Logger
}

func NewRequestHandler(logger *slog.Logger) Handler {
	return &requestHandler{
		logger: logger,
	}
}

func (h *requestHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", r.Method, "url", r.URL.Path)
	logger.Info("Handling health check request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health check passed!"))
}

func (h *requestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", r.Method, "url", r.URL.Path)
	logger.Info("Handling request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
