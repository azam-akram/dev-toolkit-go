package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type HttpServer struct {
	serv   *http.Server
	port   string
	logger *slog.Logger
}

var server Server

func NewServer(port string, logger *slog.Logger, mux *http.ServeMux) Server {
	if server != nil {
		return server
	}

	server = &HttpServer{
		port:   port,
		logger: logger,
		serv: &http.Server{
			Addr:           ":" + port,
			Handler:        mux,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    15 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
	}
	return server
}

func (s *HttpServer) Start(ctx context.Context) error {
	if err := s.serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("Error starting HTTP server", "error", err)
		return err
	}

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	if err := s.serv.Shutdown(ctx); err != nil {
		s.logger.Error("Error stopping HTTP server", "error", err)
		return err
	}
	s.logger.Info("HTTP server stopped")
	return nil
}
