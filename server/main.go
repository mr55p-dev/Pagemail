// main.go
package main

import (
	"log"
	"net/http"

	// "time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/robfig/cron/v3"

	"pagemail/server/custom_api"
)

func main() {
	app := pocketbase.New()
	c := cron.New()

	// Register the terminate handler
	app.OnTerminate().PreAdd(func(e *core.TerminateEvent) error { c.Stop(); return nil })

	// Register the app start cron handler
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		custom_api.Mailer(app)
		c.Start()
		return nil
	})

	// Register the custom routes
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			// AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
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
	if _, err := c.AddFunc("* * * * *", func() { log.Print("Cron") }); err != nil {
		log.Fatal(err)
	}

	// Start the app and cron
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
