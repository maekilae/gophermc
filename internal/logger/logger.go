package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

// --- ANSI Color Constants ---
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

// ColorConsoleHandler is a simple custom handler for colored console output
type ColorConsoleHandler struct {
	opts slog.HandlerOptions
}

func (h *ColorConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	color := reset
	switch r.Level {
	case slog.LevelDebug:
		color = blue
	case slog.LevelInfo:
		color = green
	case slog.LevelWarn:
		color = yellow
	case slog.LevelError:
		color = red
	}

	// 2. Format the time, colored level, and message
	timeStr := r.Time.Format("15:04:05")
	levelStr := fmt.Sprintf("%s%s%s", color, r.Level.String(), reset)

	// Print the base log
	fmt.Fprintf(os.Stdout, "%s %s %s", timeStr, levelStr, r.Message)

	// 3. Append attributes
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(os.Stdout, " %s=%v", a.Key, a.Value.Any())
		return true
	})

	fmt.Fprintln(os.Stdout)
	return nil
}

// Note: Implementing deep attribute/group storage for WithAttrs and WithGroup
// requires state management. For this simple wrapper, we return the handler as-is.
func (h *ColorConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ColorConsoleHandler) WithGroup(name string) slog.Handler {
	return h
}

// --- MultiHandler Implementation ---

// MultiHandler multiplexes logs to multiple slog.Handlers.
type MultiHandler struct {
	handlers []slog.Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			// Capture errors from underlying handlers
			if err := h.Handle(ctx, r); err != nil {
				errs = append(errs, err)
			}
		}
	}
	// errors.Join safely combines multiple errors (Go 1.20+)
	return errors.Join(errs...)
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

// --- Initialization and Helpers ---

// Init initializes two slog log handlers.
// It returns an error so the caller can handle initialization failures (e.g., file permissions).
func Init(logFilePrefix string) error {
	now := time.Now() // Avoid shadowing the 'time' package

	// Construct file path
	logPath := filepath.Join("logs", logFilePrefix+"_"+now.Format("20060102-150405")+".json")

	// Ensure the "logs/" directory exists before trying to open the file
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return err
	}

	// Safely handle file opening
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Use the Custom ANSI Handler for stdout
	h1 := &ColorConsoleHandler{opts: slog.HandlerOptions{Level: slog.LevelInfo}}

	// Keep the standard JSON handler for file output
	h2 := slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})

	combined := &MultiHandler{handlers: []slog.Handler{h1, h2}}
	slog.SetDefault(slog.New(combined))

	return nil
}

// Err adds a stack trace and error message to the log.
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	// debug.Stack() dynamically sizes the buffer to capture the full trace
	return slog.Group("exception",
		slog.String("message", err.Error()),
		slog.String("stacktrace", string(debug.Stack())),
	)
}
