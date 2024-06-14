package pages

import (
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
)

type ReadArticles func(userId string) ([]queries.Article, error)

func (fn ReadArticles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		return
	}
}

func NewReadArticles() ReadArticles {
	return func(userId string) ([]queries.Article, error) {
		return nil, nil
	}
}
