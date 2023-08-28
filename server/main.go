// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"pagemail/server/custom_api"
	"pagemail/server/models"
	"pagemail/server/readability"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v5"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"pagemail/server/mail"

	_ "pagemail/server/migrations"

	"github.com/robfig/cron/v3"
)

func main() {
	app := pocketbase.New()
	c := cron.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: isGoRun,
	})

	// Fetch readability config
	readerConfigDir := os.Getenv("PAGEMAIL_READABILITY_CONTEXT_DIR")
	if readerConfigDir == "" {
		panic("readability config directory not set")
	}
	readerPath, err := filepath.Abs(readerConfigDir)
	if err != nil {
		log.Panicf("Could recognise reader config dir given: %s", err)
	}
	_, err = os.Stat(readerPath)
	if err != nil {
		log.Panicf("Could not stat reader context path: %s", err)
	}
	readerConfig := models.ReaderConfig{
		NodeScript:   "main.js",
		PythonScript: "test.py",
		ContextDir:   readerConfigDir,
	}

	// Register the terminate handler
	app.OnTerminate().PreAdd(func(e *core.TerminateEvent) error { c.Stop(); return nil })

	// Register the app start cron handler
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		c.Start()
		return nil
	})

	// Register the custom routes
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/preview",
			Handler: custom_api.PreviewHandler(readerConfig),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/page/save",
			Handler: custom_api.SaveRoute(app),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				custom_api.VerifyTokenMiddleware(app),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/user/token/new",
			Handler: custom_api.NewTokenRoute(app),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/admin/mail/preview-template",
			Handler: mail.TestMailBody(readerConfig),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireAdminAuth(),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/page/readability",
			Handler: custom_api.ReadabilityHandler(app, readerConfig),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
				custom_api.ReadabilityMiddleware(app),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/page/reload",
			Handler: custom_api.ReadabilityReloadHandler(app, readerConfig),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
			},
		})

		return nil
	})

	// Register pre-write hooks
	app.OnRecordAfterCreateRequest("pages").Add(readability.PagePreviewHook(app, readerConfig))

	// Register the server cron jobs
	if _, err := c.AddFunc(
		"0 7 * * *",
		func() { mail.Mailer(app, readerConfig) },
	); err != nil {
		log.Fatal(err)
	}

	// Register commands
	app.RootCmd.AddCommand(mail.MailCommand(app, &readerConfig))
	app.RootCmd.AddCommand(readability.CrawlAll(app, readerConfig))

	// Start the app and cron
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
