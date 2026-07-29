// Package telemetry provides zero-allocation structured logging and observability primitives.
package telemetry

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Log is the global zero-allocation structured logger instance.
var Log zerolog.Logger

func init() {
	// Initialize high-performance human-readable console logger for development
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	Log = zerolog.New(output).With().Timestamp().Caller().Logger()
}
