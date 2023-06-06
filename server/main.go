// main.go
package main

import (
	"github.com/labstack/echo/v5"
	"log"
	"net/http"
	"pagemail/server/custom_api"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/robfig/cron/v3"
	"pagemail/server/mail"
)

func main() {
	app := pocketbase.New()
	c := cron.New()

	// Register the terminate handler
	app.OnTerminate().PreAdd(func(e *core.TerminateEvent) error { c.Stop(); return nil })

	// Register the app start cron handler
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		c.Start()
		return nil
	})

	// Register the custom routes
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// e.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 	AllowOrigins: []string{"*"},
		// 	// AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		// }))
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/preview",
			Handler: custom_api.Preview,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
			},
		})
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Path:    "/api/page/save",
			Handler: custom_api.SaveFactory(app),
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	// Register the server cron jobs
	if _, err := c.AddFunc(
		"0 7 * * *",
		func() { mail.Mailer(app) },
	); err != nil {
		log.Fatal(err)
	}

	// Start the app and cron
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
