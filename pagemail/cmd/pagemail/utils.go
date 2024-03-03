package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/render"
)

func MakeErrorResponse(c echo.Context, status int, err error) error {
	return render.ReturnRender(c, render.ErrorBox(err.Error()))
}

func LoadUser(c echo.Context) *db.User {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		panic(fmt.Errorf("Could not load request user"))
	}
	return user
}

func GetLoginCookie(val string) http.Cookie {
	return http.Cookie{
		Name:     auth.SESS_COOKIE,
		Value:    val,
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func IsHtmx(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func SetRedirect(c echo.Context, dest string) {
	c.Response().Header().Set("HX-Location", dest)
}
