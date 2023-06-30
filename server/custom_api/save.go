package custom_api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/pocketbase/pocketbase"
)

func saveRecord(app *pocketbase.PocketBase, user_id string, url string) error {
	collection, err := app.Dao().FindCollectionByNameOrId("pages")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)
	form.LoadData(map[string]any{
		"url":     url,
		"user_id": user_id,
	})
	return form.Submit()
}

func SaveRoute(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Fetch the page contents
		url := c.QueryParam("url")
		claims, ok := c.Get("TokenClaims").(*TokenClaims)
		if !ok {
			return c.String(http.StatusBadRequest, "Could not validate authorization token")
		}
		user_id := claims.UserID
		fmt.Printf("User ID claim: %s", user_id)

		if url == "" {
			return c.String(http.StatusBadRequest, "Must include a URL")
		}
		if user_id == "" {
			return c.String(http.StatusBadRequest, "Could not retrieve user id")
		}
		if err := saveRecord(app, user_id, url); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to store this record: %s", err))
		}
		return c.String(http.StatusOK, "Saved")
	}
}
