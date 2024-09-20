package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render"
)

// Handlers wraps all the route handlers
type Handlers struct {
	conn *sql.DB
}

// Queries gives access to the db queries object directly
func (h *Handlers) Queries() *queries.Queries {
	return queries.New(h.conn)
}

// GetIndex constructs the root page
func (s *Handlers) GetIndex(c echo.Context) error {
	return Render(c, http.StatusOK, render.Index())
}
