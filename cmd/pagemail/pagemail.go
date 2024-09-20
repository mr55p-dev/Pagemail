package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jordan-wright/email"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mr55p-dev/pagemail/assets"
	"github.com/mr55p-dev/pagemail/cmd/pagemail/urls"
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
			return strings.HasPrefix(c.Request().URL.Path, urls.Assets)
		},
		Output: os.Stdout,
	}))

	e.GET(urls.Root, srv.GetIndex)    // root
	e.GET(urls.Login, srv.GetLogin)   // login
	e.POST(urls.Login, srv.PostLogin) // login

	app := e.Group(urls.App)
	app.Use(session.Middleware(srv.store), srv.NeedsUser)
	app.GET(urls.GroupURL(urls.App, urls.App), srv.GetApp)     // app root
	app.POST(urls.GroupURL(urls.App, urls.Page), srv.PostPage) // app page

	e.StaticFS(urls.Assets, assets.FS)
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

	// connect to mail server
	mailAuth := smtp.PlainAuth(
		"",
		config.Mail.Username,
		config.Mail.Password,
		config.Mail.Host,
	)
	connPool, err := email.NewPool(
		fmt.Sprintf("%s:%d", config.Mail.Host, config.Mail.Port),
		config.Mail.PoolSize,
		mailAuth,
	)
	if err != nil {
		PanicError("Failed to start mail server", err)
	}

	testMail := &email.Email{
		From:    "Test Pagemail <mail@pagemail.io>",
		To:      []string{"Ellis <ellislunnon@gmail.com>"},
		Subject: "Test email",
		Text:    []byte("This is a test email"),
		Sender:  "mail@pagemail.io",
	}
	err = connPool.Send(testMail, time.Second*30)
	if err != nil {
		PanicError("Failed to send test email", err)
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
