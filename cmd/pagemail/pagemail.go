package main

import (
	"context"
	"crypto/rand"
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
	"github.com/mr55p-dev/pagemail/internal/preview"
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
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	logger.InfoCtx(ctx, "Starting mail job")
	client = mail.NewAwsSender(ctx, awsCfg)

	assets := getAssets(cfg.Environment)

	// Load config files
	cookieKey := MustReadFile(cfg.CookieKeyFile)
	// clientId := MustReadFile(cfg.ClientIdFile)
	// clientSecret := MustReadFile(cfg.ClientSecretFile)

	// Create the previewer and check for any "unknown" entries
	previewer := preview.New(ctx, conn)

	router, err := router.New(
		ctx,
		conn,
		assets,
		client,
		previewer,
		cookieKey,
	)
	if err != nil {
		panic(err)
	}

	// Load the mail client
	if cfg.Environment == "prd" {
		logger.Info("Starting mail client")
		go mail.MailGo(ctx, conn, router.Sender)
	} else {
		logger.Warn("Environment is not production, not starting mailGo")
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

func MustReadFile(path string) io.Reader {
	if path == "-" {
		return io.LimitReader(rand.Reader, 32)
	} else {
		logger.Debug("Using cookie key from file", "file", path)
		cookieDataFile, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		return cookieDataFile
	}
}
