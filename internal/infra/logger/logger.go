package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type FlightRecorder struct {
	logger *slog.Logger
	file   *os.File
	mu     sync.Mutex
}

var (
	instance *FlightRecorder
	once     sync.Once
)

// Init initializes the global flight recorder.
func Init(filePath string) error {
	var err error
	once.Do(func() {
		f, e := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			err = e
			return
		}

		handler := slog.NewJSONHandler(f, nil)
		logger := slog.New(handler)

		instance = &FlightRecorder{
			logger: logger,
			file:   f,
		}
	})
	return err
}

// LogThought logs a step and its details.
func LogThought(ctx context.Context, step string, details interface{}) {
	if instance == nil {
		// Fallback if not initialized, or panic depending on strictness.
		// For now, just print to stdout.
		fmt.Printf("FlightRecorder not initialized. Step: %s, Details: %v\n", step, details)
		return
	}

	instance.mu.Lock()
	defer instance.mu.Unlock()

	instance.logger.InfoContext(ctx, "thought",
		slog.String("step", step),
		slog.Any("details", details),
	)
}

// Close closes the log file.
func Close() error {
	if instance != nil && instance.file != nil {
		return instance.file.Close()
	}
	return nil
}
