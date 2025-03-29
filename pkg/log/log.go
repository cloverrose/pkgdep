package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Config is logging configuration.
type Config struct {
	Level  string // "debug", "info", "warn", "error"
	File   string // file name to write log
	Format string // "json" or "text"
}

// SetDefault initializes the default logger with the given configuration.
func SetDefault(cfg Config) (closer func() error, err error) {
	opts := &slog.HandlerOptions{
		Level: convertLogLevel(cfg.Level),
	}

	// default closer
	closer = func() error { return nil }

	// configure writer
	writer := os.Stdout
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = file
		closer = file.Close
	}

	// configure handler
	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	// set default logger
	slog.SetDefault(slog.New(handler))
	return closer, nil
}

func convertLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
