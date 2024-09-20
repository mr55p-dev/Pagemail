package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/gonk"
)

type Config struct {
	App struct {
		Host string `config:"host"`
	} `config:"app"`
}

func LogError(logger *slog.Logger, msg string, err error) {
	if err == nil {
		return
	}
	logger.Error(msg, "error", err.Error())
}

func main() {
	config := new(Config)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	yamlLoader, err := gonk.NewYamlLoader("pagemail.yaml")
	if err != nil {
		LogError(logger, "Failed to open pagemail.yaml", err)
	}
	err = gonk.LoadConfig(config, yamlLoader)
	if err != nil {
		LogError(logger, "Failed to load config", err)
	}

	srv := echo.New()
	srv.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, pagemail!")
	})

	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		if err = srv.Start(config.App.Host); err != nil {
			LogError(logger, "Failed to serve", err)
		}
		group.Done()
	}()
	group.Wait()
}
