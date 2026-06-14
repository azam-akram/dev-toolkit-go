package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	DebugWithContext(ctx context.Context, msg string, args ...any)
	InfoWithContext(ctx context.Context, msg string, args ...any)
	ErrorWithContext(ctx context.Context, msg string, args ...any)
}

type slogLogger struct {
	logger   *slog.Logger
	logLevel string
}

var (
	logger Logger
	once   sync.Once
)

func Init(level string) Logger {
	once.Do(func() {
		logger = newSlogLogger(level)
	})
	return logger
}

func Get() Logger {
	if logger == nil {
		return Init("info")
	}
	return logger
}

func newSlogLogger(level string) Logger {
	level = strings.ToLower(level)

	var handlerLevel slog.Level
	if level == "debug" {
		handlerLevel = slog.LevelDebug
	} else {
		handlerLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: handlerLevel,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// remove time and log level attribute in log item
			// it will change default log item like
			// {"time":"2022-11-08T15:28:26.000000000-05:00","level":"INFO","msg":"hello","count":3}
			// to
			// {"msg":"hello","count":3}
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	})

	return &slogLogger{
		logger:   slog.New(handler),
		logLevel: level,
	}
}

func (l *slogLogger) Debug(msg string, args ...any) {
	if l.logLevel == "debug" {
		l.logger.Debug(msg, args...)
	}
	// log nothing
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *slogLogger) DebugWithContext(ctx context.Context, msg string, args ...any) {
	if l.logLevel == "debug" {
		l.logger.DebugContext(ctx, msg, args...)
	}
	// log nothing
}

func (l *slogLogger) InfoWithContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *slogLogger) ErrorWithContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
