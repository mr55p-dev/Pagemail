package router

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func (router *Router) GetArticles(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	response.Component(render.Articles(user, []queries.Page{}, []queries.Page{}), w, r)
}
