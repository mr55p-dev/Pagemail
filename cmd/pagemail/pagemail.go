package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/router"
)

var logger = logging.NewLogger("routes")

func main() {
	ctx := context.Background()
	cfg := getConfig()
	logger := getLogger(cfg.LogLevel)

	conn := db.MustConnect(ctx, cfg.DBPath)
	defer conn.Close()

	var client mail.Sender
	if cfg.Environment == "prd" {
		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			panic(err)
		}
		logger.InfoCtx(ctx, "Starting mail job")
		client = mail.NewSesSender(ctx, awsCfg)
	} else {
		client = mail.NewNoOpSender(ctx)
	}

	assets := getAssets(cfg.Environment)
	cookieKey, err := getCookieKey(cfg.CookieKeyFile)
	if err != nil {
		panic(err)
	}

	router, err := router.New(ctx, conn, assets, client, cookieKey)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting http server", "config", cfg)
	if err := http.ListenAndServe(cfg.Host, router.Mux); err != nil {
		panic(err)
	}
}

func getConfig() *AppConfig {
	cfg := new(AppConfig)
	err := gonk.LoadConfig(
		cfg,
		gonk.FileLoader("pagemail.yaml", true),
		gonk.EnvironmentLoader("pm"),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func getLogger(level string) *logging.Logger {
	var lvl = slog.LevelInfo
	switch strings.ToUpper(level) {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "ERROR":
		lvl = slog.LevelError
	case "WARN":
		lvl = slog.LevelWarn
	case "INFO":
		lvl = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})
	logger := logging.NewLogger("main")
	logging.SetHandler(handler)
	return logger
}

func getAssets(env string) fs.FS {
	switch env {
	case "stg", "prd":
		subdir, err := fs.Sub(assets.FS, "public")
		if err != nil {
			panic(err)
		}
		return subdir
	default:
		return os.DirFS("internal/assets/public/")
	}
}

func getCookieKey(path string) (io.Reader, error) {
	if path == "-" {
		logger.Debug("Using cookie key from crpyto/rand")
		wrapper := NewBase64EncodedReader(rand.Reader)
		return wrapper, nil
	} else {
		logger.Debug("Using cookie key from file", "file", path)
		cookieDataFile, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to open cookie key file at %s: %w", path, err)
		}
		return cookieDataFile, nil
	}
}
