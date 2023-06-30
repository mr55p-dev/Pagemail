package custom_api

import (
	"net/http"
	"pagemail/server/mail"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

func MailTriggerAllFactory(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		go mail.Mailer(app)
		return c.String(http.StatusOK, "Triggered mail all in list")
	}
}
