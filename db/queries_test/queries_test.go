package queries_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/preview"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/stretchr/testify/assert"
)

var authqueries *queries.Queries
var userqueries *queries.Queries
var pagequeries *queries.Queries
var now = time.Now()
var uid = tools.GenerateNewId(5)
var pid = tools.GenerateNewId(5)

func init() {
	ctx := context.TODO()
	conn := db.MustConnect(ctx, ":memory:")
	db.MustLoadSchema(ctx, conn)
	userqueries = queries.New(conn)
	pagequeries = queries.New(conn)

	// add a test user
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		panic(err)
	}
	_, err = userqueries.WithTx(tx).CreateUser(ctx, queries.CreateUserParams{
		ID:         uid,
		Username:   "test",
		Email:      "test@mail.com",
		Subscribed: true,
	})
	if err != nil {
		panic(err)
	}
	err = authqueries.WithTx(tx).CreateLocalAuth(ctx, queries.CreateLocalAuthParams{
		UserID:       uid,
		PasswordHash: []byte("password"),
	})
	if err != nil {
		panic(err)
	}

	// add a test page
	_, err = pagequeries.WithTx(tx).CreatePage(ctx, queries.CreatePageParams{
		ID:     pid,
		UserID: uid,
		Url:    "https://example.com",
	})
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func TestReadPagesByUserBetween(t *testing.T) {
	assert := assert.New(t)
	pages, err := pagequeries.ReadPagesByUserBetween(context.TODO(), queries.ReadPagesByUserBetweenParams{
		Start:  now.Add(-time.Hour * 2),
		End:    now.Add(time.Hour * 2),
		UserID: uid,
	})
	assert.NoError(err)
	assert.Len(pages, 1)
	assert.Equal("https://example.com", pages[0].Url)
	assert.Len(pages[0].ID, 5)
}

func TestUpdatePagePreview(t *testing.T) {
	assert := assert.New(t)
	err := pagequeries.UpdatePagePreview(context.TODO(), queries.UpdatePagePreviewParams{
		Title:        sql.NullString{},
		Description:  sql.NullString{},
		ImageUrl:     sql.NullString{},
		PreviewState: preview.STATE_ERROR,
		ID:           pid,
	})
	assert.NoError(err)
}
