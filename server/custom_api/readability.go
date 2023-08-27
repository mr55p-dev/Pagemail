package custom_api

import (
	"fmt"
	"log"
	"net/http"
	"pagemail/server/models"
	"pagemail/server/readability"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/forms"
)

func ReadabilityHandler(app *pocketbase.PocketBase, readerConfig models.ReaderConfig) echo.HandlerFunc {
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

		task, err := readability.StartReaderTask(app, &page_record, readerConfig)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, task)
	}
}

func ReadabilityReloadHandler(app *pocketbase.PocketBase, readerConfig models.ReaderConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		page_id := c.QueryParam("page_id")
		if page_id == "" {
			return c.String(http.StatusBadRequest, "Must include page_id")
		}

		rawPage, err := app.Dao().FindRecordById("pages", page_id)
		if err != nil {
			return err
		}

		url := rawPage.GetString("url")
		if url == "" {
			return fmt.Errorf("Failed to fetch URL")
		}
		res, err := readability.FetchPreview(url, readerConfig)
		if err != nil {
			log.Printf("Failed to fetch preview, %s", err)
			return c.String(http.StatusInternalServerError, "Failed generating preview")
		}

		form := forms.NewRecordUpsert(app, rawPage)
		form.LoadData(res.ToMap())
		err = form.Submit()

		return c.JSON(http.StatusOK, "Done")
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

			if !record.GetBool("readability_enabled") {
				return apis.NewForbiddenError("Account is not priviledged", nil)
			}

			return next(c)
		}
	}
}
