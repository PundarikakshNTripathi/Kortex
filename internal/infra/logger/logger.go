package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// FlightRecorder is a specialized logger that records the agent's "thoughts" and actions.
// It writes structured JSON logs to a file, which can be replayed or analyzed later.
// We call it a "Flight Recorder" (like a black box on a plane) because it helps us
// understand what happened if something goes wrong.
type FlightRecorder struct {
	logger *slog.Logger // The underlying structured logger
	file   *os.File     // The file we are writing to
	mu     sync.Mutex   // Ensures we don't write from two threads at once
}

var (
	instance *FlightRecorder // Singleton instance (there can be only one)
	once     sync.Once       // Ensures Init runs only once
)

// Init initializes the global flight recorder.
// It opens the log file and sets up the JSON handler.
func Init(filePath string) error {
	var err error
	once.Do(func() {
		// Open file in append mode (don't overwrite old logs)
		f, e := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			err = e
			return
		}

		// Create a JSON handler that writes to the file
		handler := slog.NewJSONHandler(f, nil)
		logger := slog.New(handler)

		instance = &FlightRecorder{
			logger: logger,
			file:   f,
		}
	})
	return err
}

// LogThought logs a step in the agent's reasoning process.
// It takes a context (for tracing), a step name (e.g., "PLANNING"), and details.
func LogThought(ctx context.Context, step string, details interface{}) {
	if instance == nil {
		// Fallback if not initialized, just print to console
		fmt.Printf("FlightRecorder not initialized. Step: %s, Details: %v\n", step, details)
		return
	}

	instance.mu.Lock()
	defer instance.mu.Unlock()

	// Write the log entry as a JSON object
	instance.logger.InfoContext(ctx, "thought",
		slog.String("step", step),
		slog.Any("details", details),
	)
}

// Close properly closes the log file when the program exits.
func Close() error {
	if instance != nil && instance.file != nil {
		return instance.file.Close()
	}
	return nil
}
