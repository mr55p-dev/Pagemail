package readings

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/readability"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/pkg/response"
)

type ListReadings func(ctx context.Context, user *queries.User) error

func (fn ListReadings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.Get(r.Context())
	if !auth.IsAuthenticated(r.Context()) {
		response.Error(w, r, pmerror.ErrNoAuth)
		return
	}

	user := auth.GetUser(r.Context())
	err := fn(r.Context(), user)
	if err != nil {
		logger.WithError(err).InfoCtx(r.Context(), "Failed to list readings")
		response.Error(w, r, err)
		return
	}

	logger.DebugCtx(r.Context(), "Listed readings")
	response.Component(render.Readings(user), w, r)
	return
}

func NewListReadings(db *sql.DB, rbl *readability.Client) CreateReading {
	q := queries.New(db)
	return func(ctx context.Context, user *queries.User, articleID string) error {
		logger := logging.Get(ctx)

		if !user.HasReadability {
			logger.DebugCtx(ctx, "Insufficient permissions to create a reading", "user", user)
			return errNoReadabilityOnAcc
		}

		articleData, err := q.GetAllReadingInfo(ctx, user.ID)
		_ = articleData
		if err != nil {
			logger.WithError(err).ErrorCtx(ctx, "Failed to read article")
			return err
		}

		// for _, article := range articleData {
		//
		// }
		return nil
	}
}
