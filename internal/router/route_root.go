package router

import (
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func (Router) GetRoot(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	response.Component(render.Index(user), w, r)
}

func (router *Router) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Component(render.Dashboard(user, pages), w, r)
}

func (router *Router) GetAccountPage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	response.Component(render.AccountPage(user), w, r)
	return
}

type PutAccountRequest struct {
	Subscribed string `form:"email-list"`
}

func (router *Router) PutAccount(w http.ResponseWriter, r *http.Request) {
	req := request.BindRequest[PutAccountRequest](w, r)
	if req == nil {
		return
	}
	user := auth.GetUser(r.Context())
	err := router.DBClient.UpdateUserSubscription(r.Context(), dbqueries.UpdateUserSubscriptionParams{
		Subscribed: req.Subscribed == "on",
		ID:         user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (router *Router) GetShortcutToken(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	token := tools.GenerateNewShortcutToken()
	err := router.DBClient.UpdateUserShortcutToken(r.Context(), dbqueries.UpdateUserShortcutTokenParams{
		ShortcutToken: token,
		ID:            user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", token)
}
