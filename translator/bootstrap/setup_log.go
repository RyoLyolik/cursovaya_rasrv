package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// CustomHandler is a custom slog.Handler that adds a component name in brackets.
type CustomHandler struct {
	component string       // Component name (e.g., "GIN")
	handler   slog.Handler // Underlying handler for actual logging
}

// Enabled checks if the log level is enabled.
func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle formats and writes the log record.
func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	// Prepend the component name in brackets to the log message
	message := fmt.Sprintf("[%s] %s", h.component, record.Message)

	// Format additional attributes (if any)
	var attrs []string
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value.String()))
		return true
	})

	// Combine everything into a single log line
	logLine := fmt.Sprintf(
		"%s %s | %s",
		message,
		record.Time.Format("2006/01/02 - 15:04:05"),
		record.Level,
	)

	if len(attrs) > 0 {
		logLine += " | " + fmt.Sprintf("%s", attrs)
	}

	// Write the formatted log line to the underlying handler
	return h.handler.Handle(ctx, slog.Record{
		Time:    record.Time,
		Level:   record.Level,
		Message: logLine,
	})
}

// WithAttrs adds attributes to the handler.
func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		component: h.component,
		handler:   h.handler.WithAttrs(attrs),
	}
}

// WithGroup adds a group to the handler.
func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		component: h.component,
		handler:   h.handler.WithGroup(name),
	}
}

func SetupLogger(logCfg *LoggerConfig) *slog.Logger {
	var log *slog.Logger
	var level slog.Level
	customHandler := &CustomHandler{component: "APP"}

	switch logCfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	default:
		fmt.Println("Unknown log level")
		os.Exit(1)
	}

	switch logCfg.Handler {
	case "text":
		customHandler.handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	case "json":
		customHandler.handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		fmt.Println("Unknown log handler")
		os.Exit(1)
	}
	log = slog.New(customHandler.handler) // todo delete ".handler" and setup custom handler properly
	return log
}
