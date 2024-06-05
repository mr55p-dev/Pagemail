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
	pageArticles, err := q.GetPagesAndArticles(r.Context(), user.ID)
	if err != nil {
		logger.WithError(err).ErrorCtx(r.Context(), "failed to get readable pages")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}
	logger.InfoCtx(r.Context(), "Fetched article page list", "n", len(pageArticles))

	readyPages := make([]queries.Page, 0)
	pendingPages := make([]queries.Page, 0)
	unknown := make([]queries.Page, 0)
	for _, pageArticle := range pageArticles {
		article := pageArticle.Article
		if article.State == "complete" {
			readyPages = append(readyPages, pageArticle.Page)
		} else if article.State == "pending" {
			pendingPages = append(pendingPages, pageArticle.Page)
		} else {
			unknown = append(unknown, pageArticle.Page)
		}
	}
	logger.InfoCtx(
		r.Context(), "Split pages",
		"ready", len(readyPages),
		"pending", len(pendingPages),
		"not_ready", len(unknown),
	)

	response.Component(render.Articles(user, readyPages, pendingPages, unknown), w, r)
}

func (router *Router) PostReading(w http.ResponseWriter, r *http.Request) {
	logger := logger.WithRequest(r)
	user := auth.GetUser(r.Context())
	if user == nil {
		return
	}

	pageId := r.PathValue("page_id")
	logger.With("page id", pageId).Info("Requested a reading")
	q := queries.New(router.db)
	page, err := q.ReadPageById(r.Context(), pageId)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "Failed to read page")
		response.Error(w, r, pmerror.ErrNoPage)
		return
	}

	// TODO: lookup by article id and check if a reading already exists
	// if page.ReadingJobID == "" || page.ReadingJobStatus != "unknown" {
	// 	logger.InfoCtx(r.Context(),
	// 		"Tried to create reading job for page already read",
	// 	)
	// 	response.Error(w, r, pmerror.ErrReaderDuplicatePage)
	// 	return
	// }

	// TODO: fetch and extract page content... need to refactor to a service I think
	docText, err := router.Reader.Extract(r.Context(), page.Url, nil)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "failed to extract doucment text")
		response.Error(w, r, pmerror.ErrUnspecified)
		return
	}

	_ = docText
}
