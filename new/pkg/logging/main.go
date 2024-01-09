package logging

import (
	"context"
	"log/slog"
	"os"
)

type Log struct{ *slog.Logger }
type LogKey string

const (
	UserId   = "user-id"
	UserMail = "user-email"
	PageId   = "page-id"
	Error    = "error"
)

var BaseLog *slog.Logger

func init() {
	BaseLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func GetLogger(name string) Log {
	return Log{BaseLog.With("module", name)}
}

func (l *Log) Err(msg string, err error) {
	l.Error(msg, Error, err.Error())
}

func (l *Log) ErrContext(ctx context.Context, msg string, err error) {
	l.Error(msg, Error, err.Error()) 
}
