package logging

import (
	"context"
	"log/slog"
	"os"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/tools"
	"gopkg.in/yaml.v3"
)

type Log struct{ *slog.Logger }
type LogKey string

const (
	User     = "user"
	UserId   = "user-id"
	UserMail = "user-email"
	PageId   = "page-id"
	Error    = "error"
	File     = "file"
	Line     = "lineno"
	Rows     = "rows"
)

var BaseLog *slog.Logger
var Config *Cfg

type Cfg struct {
	Env      string `env:"PM_ENV" required:"true" yaml:"env" log:"environment"`
	Mode     string `env:"PM_MODE" required:"true" yaml:"mode" log:"deploy-mode"`
	DBPath   string `env:"PM_DB_PATH" required:"true" yaml:"dbpath" log:"db-path"`
	Port     string `env:"PM_PORT" required:"true" yaml:"port" log:"port"`
	TestUser string `env:"PM_TEST_USER" yaml:"testuser" log:"test-user-id"`
	LogLevel string `env:"PM_LVL" yaml:"loglevel" log:"log-level"`
}

func NewCfg() *Cfg {
	cfg := new(Cfg)
	// Parse from file
	_, err := os.Stat("pagemail.yaml")
	if err == nil {
		data, _ := os.ReadFile("pagemail.yaml")
		yaml.Unmarshal(data, cfg)
	}
	// Parse from environment
	tools.LoadFromEnv(cfg)
	err = tools.ValidateRequiredFields(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (c Cfg) LogValue() slog.Value {
	vals := tools.LogValue(&c)
	return slog.GroupValue(vals...)
}

func init() {
	Config = NewCfg()

	var handler slog.Handler
	var level slog.Level
	switch Config.LogLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	default:
	case "INFO":
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}

	if Config.Mode == "release" {
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

func (l *Log) ReqDebug(c echo.Context, msg string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.DebugContext(c.Request().Context(), msg, callerArgs...)
}
func (l *Log) ReqInfo(c echo.Context, msg string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.InfoContext(c.Request().Context(), msg, callerArgs...)
}

func (l *Log) ReqError(c echo.Context, msg string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.ErrorContext(c.Request().Context(), msg, callerArgs...)
}
func (l *Log) ReqErr(c echo.Context, msg string, err error, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{Error, err.Error(), File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.InfoContext(c.Request().Context(), msg, callerArgs...)
}
