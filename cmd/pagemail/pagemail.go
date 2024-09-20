package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/assets")
		},
		Output: os.Stdout,
	}))
	e.GET("/", srv.GetIndex)
	e.GET("/login/", srv.GetLogin)
	e.POST("/login", srv.PostLogin)
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
	cookieKey, err := os.ReadFile(config.App.CookieKeyFile)
	if err != nil {
		PanicError("Failed to read cookie key", err)
	}

	bindRoutes(server, &Handlers{
		conn:  db,
		store: sessions.NewCookieStore(cookieKey),
	})

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
