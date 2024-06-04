package pages

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
)

type CreateArticle func(ctx context.Context, pageId string) error

func (fn CreateArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		return
	}

	// page id
	pageId := r.PathValue("page_id")
	err := fn(r.Context(), pageId)
	_ = err
}

func NewCreateArticle(db *sql.DB) CreateArticle {
	q := queries.New(db)
	return func(ctx context.Context, pageId string) error {
		err := q.NewArticle(ctx, pageId)
		if err != nil {
			return err
		}

		// TODO: fetch page contents and update
		return nil
	}
}
