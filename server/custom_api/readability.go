package custom_api

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"pagemail/server/models"
	"pagemail/server/readability"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/forms"
)

func ReadabilityHandler(app *pocketbase.PocketBase, cfg *models.PMContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		page_id := c.QueryParam("page_id")
		if page_id == "" {
			return c.String(http.StatusBadRequest, "Must include page_id")
		}

		raw_page_record, err := app.Dao().FindRecordById("pages", page_id)
		if err != nil {
			return err
		}

		stat := raw_page_record.GetString("readability_status")
		if stat != string(models.ReadabilityUnknown) {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Cannot start a new synthesis job: record has status %s", stat))
		}

		page_record := models.Page{
			Id:      raw_page_record.Id,
			Url:     raw_page_record.GetString("url"),
			Created: raw_page_record.GetCreated().Time(),
		}

		readability.UpdateJobState(app, raw_page_record.Id, models.ReadabilityProcessing, nil)
		task, err := readability.StartReaderTask(app, cfg, &page_record)
		if err != nil {
			readability.UpdateJobState(app, raw_page_record.Id, models.ReadabilityFailed, nil)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Could not start reader job: %s", err))
		}

		return c.JSON(http.StatusOK, task)
	}
}

func ReadabilityReloadHandler(app *pocketbase.PocketBase, ctx *models.PMContext) echo.HandlerFunc {
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
		res, err := readability.FetchPreview(ctx, url)
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

func ReadabilityGetUrlHandler(app *pocketbase.PocketBase, ctx *models.PMContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if readability is complete
		id := c.QueryParam("page_id")
		if id == "" {
			return c.String(http.StatusBadRequest, "page_id field is missing")
		}

		record, err := app.Dao().FindRecordById("pages", id)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}

		status := record.GetString("readability_status")
		if status != string(models.ReadabilityComplete) {
			return c.String(http.StatusBadRequest, "readability processing is not complete")
		}

		// work out what the file name should be
		taskId := record.GetString("readability_task_id")
		if taskId == "" {
			return c.String(http.StatusBadRequest, "no task id found")
		}
		key := taskId + ".mp3"

		// setup s3 client
		s3Ctx := context.Background()
		client := s3.NewFromConfig(*ctx.AWS)
		signer := s3.NewPresignClient(client)
		req, err := signer.PresignGetObject(s3Ctx, &s3.GetObjectInput{
			Bucket: &ctx.S3Config.ReadabilityBucket,
			Key:    &key,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to generate presigned url: %s", err))
		}

		return c.JSON(http.StatusOK, req)
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
				return apis.NewForbiddenError("Account is not privileged", nil)
			}

			return next(c)
		}
	}
}
