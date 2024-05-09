package main

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/stretchr/testify/assert"
)

type message struct {
	address  string
	contents string
}

type MailNoOp struct {
	mail []message
}

func (m *MailNoOp) Send(ctx context.Context, addr string, contents io.Reader) error {
	dst := strings.Builder{}
	io.Copy(&dst, contents)
	cnts := dst.String()
	logger.DebugCtx(ctx, "No-op mail send", "address", addr, "contents", cnts)
	m.mail = append(m.mail, message{
		address:  addr,
		contents: cnts,
	})
	return nil
}

type MailDbReaderMock struct {
	created time.Time
}

func (MailDbReaderMock) ReadUsersWithMail(context.Context) ([]db.User, error) {
	return nil, nil
}

func newString(val string) *string {
	return &val
}

func (m *MailDbReaderMock) ReadPagesByUserBetween(context.Context, string, time.Time, time.Time) ([]db.Page, error) {

	return []db.Page{
		{
			Id:          "123",
			Url:         "https://pagemail.io",
			Title:       newString("Title 1"),
			Description: newString("Description 1"),
			Created:     &m.created,
		},
		{
			Id:          "456",
			Url:         "https://example.com",
			Title:       newString("Title 2"),
			Description: newString("Description 2"),
			Created:     &m.created,
		},
	}, nil
}

func TestSendMailToUser(t *testing.T) {
	assert := assert.New(t)
	testUser := db.User{
		Id:         "123",
		Username:   "user",
		Email:      "user@mail.com",
		Name:       "User",
		Subscribed: true,
	}

	created := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	now := time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)
	sender := MailNoOp{}
	dbReader := &MailDbReaderMock{
		created: created,
	}
	err := sendMailToUser(context.TODO(), &testUser, &sender, dbReader, now)
	assert.NoError(err)
	assert.Len(sender.mail, 1)
}
