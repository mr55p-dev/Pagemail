package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mr55p-dev/pagemail/assets"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/render"
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
	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, "/assets")
			},
			Output: os.Stdout,
		}),
	)
	e.Pre(middleware.RemoveTrailingSlash())
	authMiddlewares := []echo.MiddlewareFunc{session.Middleware(srv.store), srv.NeedsUser}

	e.GET("/", srv.GetPage(render.Index))
	e.GET("/login", srv.GetPage(render.Login))
	e.GET("/signup", srv.GetPage(render.Signup))
	e.POST("/login", srv.PostLogin)
	s.POST("/signup", srv.PostSignup)
	e.GET("/logout", srv.GetLogout, authMiddlewares...)

	app := e.Group("/app", authMiddlewares...)
	app.Use(session.Middleware(srv.store), srv.NeedsUser)

	app.GET("", srv.GetApp)
	app.POST("/page", srv.PostPage)
	app.DELETE("/page/:id", srv.DeletePage)

	e.StaticFS("/assets", assets.FS)
}

func concatHostPort(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
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

	// create the mail pool
	mailPool, err := mail.NewPool(
		config.Mail.Username,
		config.Mail.Password,
		config.Mail.Host,
		config.Mail.Port,
		config.Mail.PoolSize,
	)
	if err != nil {
		PanicError("Failed to open mail pool", err)
	}
	mailTimeout := time.Minute
	mailInterval := time.Minute * 30
	mailer := mail.New(ctx, db, mailPool, mailTimeout)
	go func() {
		for {
			timer := time.NewTimer(mailInterval)
			select {
			case now := <-timer.C:
				count, err := mailer.RunScheduledSend(ctx, now)
				logger.InfoContext(ctx, "Finished mail job", "count", count, "errors", err)
			case <-ctx.Done():
				timer.Stop()
				return
			}

		}
	}()

	// bind everything together
	bindRoutes(server, &Handlers{
		conn:  db,
		store: sessions.NewCookieStore(cookieKey),
		mail:  mailPool,
	})

	// start the server
	go func() {
		defer appCancel()
		LogError("Failed to serve", server.Start(concatHostPort(config.App.Host, config.App.Port)))
	}()

	// wait for an exit condition
	select {
	case <-ctx.Done():
		logger.Info("Closing application", "reason", "app context canclled")
	case <-interruptChan:
		logger.Info("Interrupt received, shutting down")
	}
	return
}
