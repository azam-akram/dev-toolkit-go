package circuitbreaker

import (
	"errors"
	"time"
)

type Config struct {
	FailureThreshold int
	FailureWindow    time.Duration
	OpenDuration     time.Duration
}

func (c Config) Validate() error {
	if c.FailureThreshold <= 0 {
		return errors.New("failure threshold must be > 0")
	}

	if c.FailureWindow <= 0 {
		return errors.New("failure window must be > 0")
	}

	if c.OpenDuration <= 0 {
		return errors.New("open duration must be > 0")
	}

	return nil
}
