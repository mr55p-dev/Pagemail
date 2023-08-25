package custom_api

import (
	"net/http"
	"pagemail/server/mail"
	"pagemail/server/readability"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

func MailTriggerAllFactory(app *pocketbase.PocketBase, cfg readability.ReaderConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		go mail.Mailer(app, cfg)
		return c.String(http.StatusOK, "Triggered mail all in list")
	}
}
