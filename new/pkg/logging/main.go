package logging

import (
	"context"
	"log/slog"
	"os"
	"runtime"
)

type Log struct{ *slog.Logger }
type LogKey string

const (
	UserId   = "user-id"
	UserMail = "user-email"
	PageId   = "page-id"
	Error    = "error"
	File     = "file"
	Line     = "lineno"
)

var BaseLog *slog.Logger

func init() {
	env := os.Getenv("PM_ENV")
	lvl := os.Getenv("PM_LVL")
	var handler slog.Handler
	var level slog.Level
	switch lvl {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	default:
	case "INFO":
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}

	if env == "prd" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	BaseLog = slog.New(handler)
}

func GetLogger(name string) Log {
	return Log{BaseLog.With("module", name)}
}

func (l *Log) Err(msg string, err error, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{Error, err.Error(), File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.Error(msg, callerArgs...)
}

func (l *Log) ErrContext(ctx context.Context, msg string, err error, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{Error, err.Error(), File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.ErrorContext(ctx, msg, callerArgs...)
}
