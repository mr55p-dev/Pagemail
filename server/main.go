// main.go
package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type UrlData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func fetch_url(url string, data *UrlData) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	title := doc.Find("meta[property='og:title']").AttrOr("content", "")
	description := doc.Find("meta[property='og:description']").AttrOr("content", "")

	//Fallbacks
	if title == "" {
		title = doc.Find("title").Text()
	}
	if description == "" {
		doc.Find("meta[name='description']").Each(func(i int, sel *goquery.Selection) {
			if desc, exists := sel.Attr("content"); exists && desc != "" {
				description = desc
			}
		})
	}

	data.Title = title
	data.Description = description
	return nil
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/preview",
			Handler: func(c echo.Context) error {
				// Fetch the page contents
				uri := c.QueryParam("target")
				if uri == "" {
					return c.String(http.StatusBadRequest, "Must include a URL")
				}
				data := new(UrlData)
				if err := fetch_url(uri, data); err != nil {
					return c.String(http.StatusServiceUnavailable, "Failed to fetch the external resource")
				}
				return c.JSON(http.StatusOK, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				apis.RequireRecordAuth("users"),
			},
		})

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/page/save",
			Handler: func(c echo.Context) error {
				// Fetch the page contents
				uri := c.FormValue("target")
				user_id := c.FormValue("uid")
				if uri == "" {
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
				record.Set("url", uri)
				record.Set("user_id", user_id)

				if err := app.Dao().SaveRecord(record); err != nil {
					return c.String(http.StatusBadRequest, "Failed to store this record")
				}
				return c.String(http.StatusOK, "Saved")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				// apis.RequireRecordAuth("users"),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
