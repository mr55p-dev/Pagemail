package queries_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/db/queries/queryauth"
	"github.com/mr55p-dev/pagemail/db/queries/querypages"
	"github.com/mr55p-dev/pagemail/db/queries/queryusers"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/stretchr/testify/assert"
)

var authqueries *queryauth.Queries
var userqueries *queryusers.Queries
var pagequeries *querypages.Queries
var now = time.Now()
var uid = tools.GenerateNewId(5)
var pid = tools.GenerateNewId(5)

func init() {
	ctx := context.TODO()
	conn := db.MustConnect(ctx, ":memory:")
	db.MustLoadSchema(ctx, conn)
	userqueries = queryusers.New(conn)
	pagequeries = querypages.New(conn)

	// add a test user
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		panic(err)
	}
	err = userqueries.WithTx(tx).CreateUser(ctx, queryusers.CreateUserParams{
		ID:         uid,
		Username:   "test",
		Email:      "test@mail.com",
		Subscribed: true,
	})
	if err != nil {
		panic(err)
	}
	err = authqueries.WithTx(tx).CreateLocalAuth(ctx, queryauth.CreateLocalAuthParams{
		UserID:       uid,
		PasswordHash: []byte("password"),
	})
	if err != nil {
		panic(err)
	}

	// add a test page
	err = pagequeries.WithTx(tx).CreatePage(ctx, querypages.CreatePageParams{
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
	pages, err := pagequeries.ReadPagesByUserBetween(context.TODO(), querypages.ReadPagesByUserBetweenParams{
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
	err := pagequeries.UpdatePagePreview(context.TODO(), querypages.UpdatePagePreviewParams{
		Title:        sql.NullString{},
		Description:  sql.NullString{},
		ImageUrl:     sql.NullString{},
		PreviewState: "failed",
		ID:           pid,
	})
	assert.NoError(err)
}
