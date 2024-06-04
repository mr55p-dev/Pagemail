package readings

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/readability"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

var (
	errArticleNotFound = &pmerror.PMError{
		Message: "Article not found",
		Status:  http.StatusNotFound,
	}
	errNoReadabilityOnAcc = &pmerror.PMError{
		Message: "You do not have the correct access rights to create readings",
		Status:  http.StatusForbidden,
	}
)

type CreateReading func(ctx context.Context, user *queries.User, articleID string) error

func (fn CreateReading) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.Get(r.Context())
	if !auth.IsAuthenticated(r.Context()) {
		response.Error(w, r, pmerror.ErrNoAuth)
		return
	}

	user := auth.GetUser(r.Context())
	articleID := r.PathValue("article_id")
	err := fn(r.Context(), user, articleID)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "Failed to create reading")
		response.Error(w, r, err)
		return
	}

	logger.InfoCtx(r.Context(), "Created reading job")
	return
}

func NewCreateReading(db *sql.DB, rbl *readability.Client) CreateReading {
	q := queries.New(db)
	return func(ctx context.Context, user *queries.User, articleID string) error {
		logger := logging.Get(ctx)

		if !user.HasReadability {
			logger.DebugCtx(ctx, "Insufficient permissions to create a reading", "user", user)
			return errNoReadabilityOnAcc
		}

		article, err := q.GetArticle(ctx, articleID)
		if err != nil {
			logger.WithError(err).ErrorCtx(ctx, "Failed to read article")
			return err
		}

		if article.UserID != user.ID {
			logger.DebugCtx(ctx, "access violation for article",
				"article_id", article.ID,
				"owner_id", article.UserID,
				"user_id", user.ID,
			)
			return errArticleNotFound
		}

		logger.InfoCtx(ctx, "Requesting new reading for article", "article_id", articleID)
		res, err := rbl.Synthesize(ctx, bytes.NewReader(article.Content))
		if err != nil {
			logger.WithError(err).ErrorCtx(ctx, "Failed to create reading job")
			return err
		}
		if res.JobId != "" {
			_, err = q.NewReading(ctx, queries.NewReadingParams{
				ID:        tools.NewReadingId(),
				UserID:    user.ID,
				ArticleID: articleID,
				JobID:     res.JobId,
				State:     res.Status,
			})
		}
		if len(res.Errors) > 0 {
			errs := make([]error, len(res.Errors))
			for i, e := range res.Errors {
				errs[i] = fmt.Errorf("Readability client error: %s (%s)", e.Message, e.Detail)
			}
			return errors.Join(errs...)
		}
		if err != nil {
			return err
		}

		return nil
	}
}
