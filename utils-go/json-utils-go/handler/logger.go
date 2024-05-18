package handler

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once
var log *zap.Logger

type adaptedLogger struct {
	zapLogger *zap.Logger
}

// GetLogger - get/create logger. Info level by default.
func GetLogger() Logger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		log, _ = config.Build()
	})

	return &adaptedLogger{log}
}

type Logger interface {
	Print(args ...interface{})
	PrintKeyValue(msg string, keysAndValues ...interface{})
}

func (l *adaptedLogger) Print(args ...interface{}) {
	log.Sugar().Info(args...)
}

func (l *adaptedLogger) PrintKeyValue(msg string, keysAndValues ...interface{}) {
	log.Sugar().Infow(msg, keysAndValues...)
}
