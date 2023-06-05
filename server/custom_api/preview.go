package custom_api

import (
	"github.com/labstack/echo/v5"
	"github.com/PuerkitoBio/goquery"
	"net/http"
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

func Preview(c echo.Context) error {
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
}
