package custom_api

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"net/http"

	"github.com/pocketbase/pocketbase"
)

func SaveFactory(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Fetch the page contents
		url := c.FormValue("url")
		user_id := c.FormValue("user_id")
		if url == "" {
			return c.String(http.StatusBadRequest, "Must include a URL")
		}
		if user_id == "" {
			return c.String(http.StatusBadRequest, "Must include a user id")
		}
		collection, err := app.Dao().FindCollectionByNameOrId("pages")
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}
		record := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, record)
		form.LoadRequest(c.Request(), "")
		if err := form.Submit(); err != nil {
			return c.String(http.StatusBadRequest, "Failed to store this record")
		}
		return c.String(http.StatusOK, "Saved")
	}
}
