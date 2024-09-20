package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/render"
)

type Config struct {
	App struct {
		Host string `config:"host"`
	} `config:"app"`
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func LogError(logger *slog.Logger, msg string, err error) {
	if err == nil {
		return
	}
	logger.Error(msg, "error", err.Error())
}

type Server struct {
	*echo.Echo
}

func main() {
	config := new(Config)

	appContext, appCancel := context.WithCancel(context.Background())
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)

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

	srv := &Server{echo.New()}
	bindRoutes(srv)

	go func() {
		defer appCancel()
		if err = srv.Start(config.App.Host); err != nil {
			LogError(logger, "Failed to serve", err)
		}
	}()
	select {
	case <-appContext.Done():
		logger.Info("Closing application", "reason", "app context canclled")
	case <-interruptChan:
		logger.Info("Interrupt received, shutting down")
	}
	return
}

func bindRoutes(server *Server) {
	server.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, render.Index())
	})
}
