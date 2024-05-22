package router

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

func (router *Router) GetArticles(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())
	q := queries.New(router.db)
	logger.InfoCtx(r.Context(), "Getting articles for user")
	readablePages, err := q.ReadPagesByReadable(r.Context(), queries.ReadPagesByReadableParams{
		Readable: true,
		UserID:   user.ID,
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "failed to get readable pages")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}
	logger.InfoCtx(r.Context(), "Fetched article page list", "n", len(readablePages))

	readyPages := make([]queries.Page, 0)
	pendingPages := make([]queries.Page, 0)
	notReadyPages := make([]queries.Page, 0)
	for _, page := range readablePages {
		if page.ReadingJobStatus == "complete" {
			readyPages = append(readyPages, page)
		} else if page.ReadingJobStatus == "pending" {
			pendingPages = append(pendingPages, page)
		} else {
			notReadyPages = append(notReadyPages, page)
		}
	}
	logger.InfoCtx(
		r.Context(), "Split pages",
		"ready", len(readyPages),
		"pending", len(pendingPages),
		"not_ready", len(notReadyPages),
	)

	response.Component(render.Articles(user, readyPages, pendingPages, notReadyPages), w, r)
}