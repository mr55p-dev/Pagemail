package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/pkg/auth"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/middlewares"
	"github.com/mr55p-dev/pagemail/pkg/render"
)

type Router struct {
	DBClient   db.Client
	Authorizer auth.AbsAuthorizer
}

type DataIndex struct {
	IsUser bool
}

type DataPages struct {
	UserId string
	Pages  []db.Page
}

func (Router) GetRoot(c echo.Context) error {
	return render.RenderTempate("index", c.Response(), &DataIndex{false})
}

func (Router) GetLogin(c echo.Context) error {
	return render.RenderTempate("login", c.Response(), nil)
}

func (s *Router) PostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	user, err := s.DBClient.ReadUserByEmail(email)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if !s.Authorizer.ValidateUser(email, password, user) {
		return c.String(http.StatusBadRequest, "Invalid username or password")
	}

	sess := s.Authorizer.GetToken(user.Id)

	c.Response().Header().Set("Location", fmt.Sprintf("/%s/pages", user.Id))
	c.SetCookie(&http.Cookie{
		Name:  auth.SESS_COOKIE,
		Value: sess,
		Path:  "/",
	})
	return c.NoContent(http.StatusSeeOther)
}

func (Router) GetSignup(c echo.Context) error {
	return render.RenderTempate("signup", c.Response(), nil)
}

func (s *Router) PostSignup(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordRepeat := c.FormValue("password-repeat")
	if password != passwordRepeat {
		return c.String(http.StatusBadRequest, "Passwords do not match")
	}

	token, err := s.Authorizer.SignupNewUser(email, password, username)
	if err != nil {
		return c.String(http.StatusBadRequest, "Something went wrong")
	}

	c.Response().Header().Set("Location", "/pages")
	c.SetCookie(&http.Cookie{
		Name:  auth.SESS_COOKIE,
		Value: token,
		Path:  "/",
	})
	return c.NoContent(http.StatusSeeOther)
}

func (Router) GetLogout(c echo.Context) error {
	c.Response().Header().Set("Location", "/login")
	c.SetCookie(&http.Cookie{
		Name:   auth.SESS_COOKIE,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return c.NoContent(http.StatusSeeOther)
}

func (r *Router) GetPages(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	pages, err := r.DBClient.ReadPagesByUserId(user.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return render.RenderTempate("pages", c.Response(), &DataPages{
		Pages:  pages,
		UserId: user.Id,
	})
}

func (r *Router) DeletePages(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	n, err := r.DBClient.DeletePagesByUserId(user.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("Deleted %d records", n))

}

func (r *Router) PostPage(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok || user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	url := c.FormValue("url")
	page := db.NewPage(user.Id, url)

	if err := r.DBClient.CreatePage(page); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (r *Router) ListenPages(c echo.Context) error {
	user, ok := c.Get("user").(*db.User)
	if !ok || user == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// create listener to updates
	events := make(chan db.Event[db.Page])

	r.DBClient.AddPageListener(user.Id, events)

	c.Response().WriteHeader(http.StatusOK)
	c.Response().Header().Add("Content-Type", "text/event-stream")
	c.Response().Flush()

	defer r.DBClient.RemovePageListener(user.Id)
	for event := range events {
		fmt.Println("Event", event)
		c.Response().Write([]byte("<p>Message sent</p>\n\n"))
		c.Response().Flush()
	}

	return nil
}

func (r *Router) TestUpdate(c echo.Context) error {
	user, _ := c.Get("user").(*db.User)
	page := db.NewPage(user.Id, "https://example.com")
	page.Id = "123456"
	err := r.DBClient.UpsertPage(page)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(200, "Updated")
}

func (Router) TestTmpl(c echo.Context) error {
	cmp := render.Hello("ellis")
	return cmp.Render(c.Request().Context(), c.Response())
}

func main() {
	// Logging
	e := echo.New()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	log := slog.New(handler)
	dbLogger := log.With("module", "db")
	authLogger := log.With("module", "auth")
	protectionLogger := log.With("module", "protection middleware")
	reqLogger := log.With("module", "request log middleware")

	dbClient := db.NewDriver(dbLogger)
	defer dbClient.Close()

	authClient := auth.NewAuthorizer(dbClient, authLogger)

	s := &Router{
		DBClient:   dbClient,
		Authorizer: authClient,
	}

	e.Use(middlewares.TraceMiddleware)
	protected := middlewares.GetProtectedMiddleware(protectionLogger, authClient, dbClient)
	e.Use(middlewares.GetLoggingMiddleware(reqLogger))

	e.Static("/assets", "public")
	e.GET("/", s.GetRoot)

	e.GET("/login", s.GetLogin)
	e.POST("/login", s.PostLogin)
	e.GET("/logout", s.GetLogout, protected)

	e.GET("/signup", s.GetSignup)
	e.POST("/signup", s.PostSignup)

	e.GET("/:id/pages", s.GetPages, protected)
	e.DELETE("/:id/pages", s.DeletePages, protected)
	e.POST("/:id/page", s.PostPage, protected)
	e.GET("/:id/pages/listen", s.ListenPages, protected)
	e.GET("/test", s.TestUpdate, protected)

	if err := e.Start(":8080"); err != nil {
		log.Error(err.Error())
	}
}
