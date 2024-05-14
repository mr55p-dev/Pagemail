package router

import (
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

type AccountData struct {
	Subscribed string `form:"email-list"`
}

func (Router) GetRoot(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	componentRender(render.Index(user), w, r)
}

func (router *Router) GetDashboard(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	pages, err := router.DBClient.ReadPagesByUserId(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	componentRender(render.Dashboard(user, pages), w, r)
}

func (router *Router) GetAccountPage(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	componentRender(render.AccountPage(user), w, r)
	return
}

func (router *Router) PutAccount(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
	form := new(AccountData)
	err := router.DBClient.UpdateUserSubscription(r.Context(), dbqueries.UpdateUserSubscriptionParams{
		Subscribed: form.Subscribed == "on",
		ID:         user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (router *Router) GetShortcutToken(w http.ResponseWriter, r *http.Request) {
	user := dbqueries.GetUser(r.Context())
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
