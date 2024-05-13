package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/router"
)

var logger = logging.NewLogger("routes")

func main() {
	ctx := context.Background()
	cfg := getConfig()
	awsCfg := getAwsConfig(ctx)
	logger := getLogger(cfg)

	router, err := router.New(ctx, cfg, awsCfg)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting http server", "config", cfg)
	if err := http.ListenAndServe(cfg.Host, router.Mux); err != nil {
		panic(err)
	}
}

func getConfig() *router.AppConfig {
	cfg := new(router.AppConfig)
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

func getAwsConfig(ctx context.Context) aws.Config {
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	return awsCfg
}

func getLogger(cfg *router.AppConfig) *logging.Logger {
	var lvl = slog.LevelInfo
	switch strings.ToUpper(cfg.LogLevel) {
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
