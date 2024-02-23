package render

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func ReturnRender(c echo.Context, component templ.Component) error {
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response().Writer)
}
