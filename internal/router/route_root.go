package router

import (
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func (Router) GetRoot(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	response.Component(render.Index(user), w, r)
}

func (router *Router) GetAccountPage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	response.Component(render.AccountPage(user), w, r)
	return
}

type PutAccountRequest struct {
	Subscribed string `form:"subscribed"`
}

func (router *Router) PutAccount(w http.ResponseWriter, r *http.Request) {
	req := request.BindRequest[PutAccountRequest](w, r)
	if req == nil {
		return
	}
	user := auth.GetUser(r.Context())
	err := queries.New(router.db).UpdateUserSubscription(r.Context(), queries.UpdateUserSubscriptionParams{
		Subscribed: req.Subscribed == "on",
		ID:         user.ID,
	})
	if err != nil {
		response.Error(w, r, pmerror.NewInternalError("Failed to update account"))
		return
	}
	response.Success("Updated account", w, r)
}

func (router *Router) GetShortcutToken(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	token, tokenHash := auth.NewShortcutToken()
	err := queries.New(router.db).UpdateShortcutToken(r.Context(), queries.UpdateShortcutTokenParams{
		UserID:     user.ID,
		Credential: tokenHash,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(fmt.Sprintf("Generated new shortcut token: %s", token), w, r)
}
