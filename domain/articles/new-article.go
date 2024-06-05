package pages

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/request"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

type CreateArticle func(ctx context.Context, user *queries.User, pageID string) error

type CreateArticleParams struct {
	pageID string `query:"page_id"`
}

func (fn CreateArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.Get(r.Context())
	params := request.BindRequest[CreateArticleParams](w, r)
	if params == nil {
		return
	}
	user := auth.GetUser(r.Context())
	if user == nil {
		return
	}

	err := fn(r.Context(), user, params.pageID)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "Failed to create new article")
		response.Error(w, r, err)
		return
	}
	logger.InfoCtx(r.Context(), "Created new article")
	response.Success("Creaed article", w, r)
	return
}

func NewCreateArticle(db *sql.DB) CreateArticle {
	q := queries.New(db)
	return func(ctx context.Context, user *queries.User, pageID string) error {
		err := q.NewArticle(ctx, queries.NewArticleParams{
			ID:     tools.NewArticleId(),
			UserID: user.ID,
			PageID: pageID,
		})
		if err != nil {
			return err
		}

		// TODO: fetch page contents and update
		return nil
	}
}
