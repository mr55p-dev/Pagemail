package logging

import (
	"context"
	"io"
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Err(msg string, err error) {
	l.Error(msg, "error", err.Error())
}

func (l *Logger) Errc(ctx context.Context, msg string, err error) {
	l.ErrorContext(ctx, msg, "error", err.Error())
}

func NewVoid() *Logger {
	return &Logger{
		Logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
	}
}

func New(log *slog.Logger) *Logger {
	return &Logger{
		Logger: log,
	}
}
