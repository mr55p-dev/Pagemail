package custom_api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"pagemail/server/models"
	"pagemail/server/readability"
)

func PreviewHandler(ctx *models.PMContext) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Fetch the page contents
		uri := c.QueryParam("target")
		if uri == "" {
			return c.String(http.StatusBadRequest, "Must include a URL")
		}
		data, err := readability.FetchPreview(ctx, uri)
		if err != nil {
			return c.String(http.StatusServiceUnavailable, "Failed to fetch the external resource")
		}
		c.Response().Header().Set("Cache-Control", "private, max-age=432000")
		return c.JSON(http.StatusOK, data)
	}
}
