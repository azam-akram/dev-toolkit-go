package utils

import (
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/logger"
	"encoding/json"
)

func StringToStruct(s string, o any) error {
	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		logger.Get().Error("Failed to convert string into struct", "error", err)
	}

	return err
}
