// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"pagemail/server/custom_api"
	"pagemail/server/models"
	"pagemail/server/readability"
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

	// Load the configuration
	ctx := new(models.PMContext)
	ctx.LoadFromEnv()

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
			Handler: custom_api.PreviewHandler(ctx),
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
			Handler: mail.TestMailBody(),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireAdminAuth(),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/page/readability",
			Handler: custom_api.ReadabilityHandler(app, ctx),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
				custom_api.ReadabilityMiddleware(app),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/page/reload",
			Handler: custom_api.ReadabilityReloadHandler(app, ctx),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
			},
		})

		return nil
	})

	// Register pre-write hooks
	app.OnRecordAfterCreateRequest("pages").Add(readability.PagePreviewHook(app, ctx))

	// Register the server cron jobs
	if _, err := c.AddFunc(
		"0 7 * * *",
		func() { mail.Mailer(app, ctx) },
	); err != nil {
		log.Fatal(err)
	}

	// Register commands
	app.RootCmd.AddCommand(mail.MailCommand(app, ctx))
	app.RootCmd.AddCommand(readability.CrawlAll(app, ctx))

	// Start the app and cron
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
