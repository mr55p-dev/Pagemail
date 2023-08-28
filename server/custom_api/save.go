package custom_api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"

	"github.com/pocketbase/pocketbase"
)

func saveRecord(app *pocketbase.PocketBase, collection *models.Collection, user_id string, url string) (*models.Record, error) {
	record := models.NewRecord(collection)
	record.Load(map[string]any{
		"url":     url,
		"user_id": user_id,
	})
	return record, nil
}

func SaveRoute(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		collection, err := app.Dao().FindCollectionByNameOrId("pages")
		if err != nil {
			log.Panic(err)
		}

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
		
		record, err := saveRecord(app, collection, user_id, url)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to store this record: %s", err))
		}

		form := forms.NewRecordUpsert(app, record)
		err = form.Submit()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to save record")
		}

		go func() {
			err = app.OnRecordAfterCreateRequest("pages").Trigger(&core.RecordCreateEvent{
				HttpContext: c,
				Record:      record,
				BaseCollectionEvent: core.BaseCollectionEvent{
					Collection: collection,
				},
			})
		}()

		return c.String(http.StatusOK, "Saved")
	}
}
