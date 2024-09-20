package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/assets"
)

// Global logger instance
var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
}))

// Global config instance
var config = MustLoadConfig()

// bindRoutes attaches route handlers to endpoints
func bindRoutes(e *echo.Echo, srv *Handlers) {
	e.GET("/", srv.GetIndex)
	e.StaticFS("/assets", assets.FS)
}

func main() {
	// configure interrupt handling
	ctx, appCancel := context.WithCancel(context.Background())
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)

	// connect to DB
	db, err := openDB(ctx, config.DB.Path)
	if err != nil {
		PanicError("Failed to open db connection", err)
	}

	// create the routes
	server := echo.New()
	handlers := &Handlers{conn: db}
	bindRoutes(server, handlers)

	// start the server
	go func() {
		defer appCancel()
		if err := server.Start(config.App.Host); err != nil {
			LogError("Failed to serve", err)
		}
	}()
	select {
	case <-ctx.Done():
		logger.Info("Closing application", "reason", "app context canclled")
	case <-interruptChan:
		logger.Info("Interrupt received, shutting down")
	}
	return
}
