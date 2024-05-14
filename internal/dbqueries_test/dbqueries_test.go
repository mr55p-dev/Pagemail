package dbqueries_test

import (
	"context"
	"testing"
	"time"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/stretchr/testify/assert"
)

var queries *dbqueries.Queries
var now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
var uid = tools.GenerateNewId(5)

func init() {
	ctx := context.TODO()
	conn := dbqueries.MustGetDB(ctx, ":memory:")
	queries = dbqueries.New(conn)
	err := dbqueries.LoadTables(ctx, conn)
	if err != nil {
		panic(err)
	}

	// add a test user
	err = queries.CreateUser(ctx, dbqueries.CreateUserParams{
		ID:             uid,
		Username:       "test",
		Email:          "test@mail.com",
		Password:       []byte("password"),
		Subscribed:     true,
		ShortcutToken:  tools.GenerateNewShortcutToken(),
		HasReadability: false,
		Created:        now,
		Updated:        now,
	})
	if err != nil {
		panic(err)
	}

	// add a test page
	err = queries.CreatePage(ctx, dbqueries.CreatePageParams{
		ID:      tools.GenerateNewId(5),
		UserID:  uid,
		Url:     "https://example.com",
		Created: now.Add(-time.Hour),
		Updated: now.Add(-time.Hour),
	})
	if err != nil {
		panic(err)
	}
}

func TestReadPagesByUserBetween(t *testing.T) {
	assert := assert.New(t)
	pages, err := queries.ReadPagesByUserBetween(context.TODO(), dbqueries.ReadPagesByUserBetweenParams{
		Start:  now.Add(-time.Hour * 2),
		End:    now,
		UserID: uid,
	})
	assert.NoError(err)
	assert.Len(pages, 1)
	assert.Equal("https://example.com", pages[0].Url)
	assert.Len(pages[0].ID, 5)
}
