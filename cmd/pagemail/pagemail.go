package main

import (
	"context"
	"crypto/rand"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/internal/assets"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/readability"
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
	cookieKey := MustReadFile(cfg.CookieKeyFile)

	reader := readability.New(ctx, &url.URL{
		Scheme: cfg.Readability.Scheme,
		Host:   cfg.Readability.Host,
	})
	if err != nil {
		panic(err)
	}
	previewer := preview.New(ctx, conn, reader)

	router, err := router.New(
		ctx,
		conn,
		assets,
		client,
		previewer,
		cookieKey,
		cfg.GoogleClientId,
		cfg.External.Host,
		cfg.External.Scheme,
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

	sources := make([]gonk.Loader, 0)
	yamlSource, err := gonk.NewYamlLoader("pagemail.yaml")
	if err == nil {
		sources = append(sources, yamlSource)
	} else {
		logger.WithError(err).Warn("Could not load local pagemail.yaml file")
	}
	sources = append(sources, gonk.EnvLoader("pm"))
	err = gonk.LoadConfig(cfg, sources...)
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
	logging.Level.Set(lvl)
	logger := logging.NewLogger("main")
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
