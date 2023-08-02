package custom_api

import (
	"net/http"
	"pagemail/server/models"
	"pagemail/server/readability"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
)

func ReadabilityHandler(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		page_id := c.QueryParam("page_id")
		if page_id == "" {
			return c.String(http.StatusBadRequest, "Must include page_id")
		}

		raw_page_record, err := app.Dao().FindRecordById("pages", page_id)
		if err != nil {
			return err
		}

		page_record := models.Page{
			Url:     raw_page_record.GetString("url"),
			Created: raw_page_record.GetCreated().Time(),
		}

		task, err := readability.StartReaderTask(app, &page_record)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, task)
	}
}


func ReadabilityMiddleware(app *pocketbase.PocketBase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, err := app.Dao().FindAuthRecordByToken(
				c.Request().Header.Get("Authorization"),
				app.Settings().RecordAuthToken.Secret,
			)
			if err != nil {
				return apis.NewUnauthorizedError("Failed to find user", err)
			}

			if !record.GetBool("isReadabilityAllowed") {
				return apis.NewForbiddenError("Account is not priviledged", nil)
			}

			return next(c)
		}
	}
}
