// logger/logger.go

package logger

import (
	"os"

	"log/slog"
)

var logger *slog.Logger

// InitializeLogger sets up the global logger.
// It should be called once, typically at application startup.
func InitializeLogger() *slog.Logger {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	return logger
}

// SetLogger allows setting a custom logger, useful for testing.
func SetLogger(l *slog.Logger) {
	logger = l
}

// GetLogger returns the current logger.
func GetLogger() *slog.Logger {
	if logger != nil {
		return logger
	}
	return slog.Default()
}
