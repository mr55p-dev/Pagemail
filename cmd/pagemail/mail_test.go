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

type mailSenderNoOp struct {
	mail []message
}

func (m *mailSenderNoOp) Send(ctx context.Context, addr string, contents io.Reader) error {
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

func (m *mailSenderNoOp) Reset() {
	m.mail = make([]message, 0)
}

type mailDbReaderNoOp struct {
	created  time.Time
	user     db.User
	numPages int
	numUsers int
}

func (m *mailDbReaderNoOp) ReadUsersWithMail(context.Context) ([]db.User, error) {
	users := make([]db.User, 0)
	for i := 0; i < m.numUsers; i++ {
		users = append(users, testUser)
	}
	return users, nil
}

func (m *mailDbReaderNoOp) ReadPagesByUserBetween(context.Context, string, time.Time, time.Time) ([]db.Page, error) {
	testPages := []db.Page{
		{
			Id:      "123",
			UserId:  m.user.Id,
			Url:     "https://pagemail.io",
			Created: &m.created,
		},
		{
			Id:      "123",
			UserId:  m.user.Id,
			Url:     "https://example.com",
			Created: &m.created,
		},
	}
	if m.numPages == 0 {
		return []db.Page{}, nil
	}
	return testPages[:m.numPages], nil
}

var (
	created  = time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	now      = time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)
	testUser = db.User{
		Id:         "123",
		Username:   "user",
		Email:      "user@mail.com",
		Name:       "User",
		Subscribed: true,
	}

	mailSender = &mailSenderNoOp{}
	dbReader   = &mailDbReaderNoOp{
		created:  created,
		numPages: 1,
		numUsers: 1,
	}
)

func TestSendMailToUser(t *testing.T) {
	defer mailSender.Reset()
	dbReader.numUsers = 1
	assert := assert.New(t)
	err := SendMailToUser(context.TODO(), &testUser, dbReader, mailSender, now)
	assert.NoError(err)
	assert.Len(mailSender.mail, 1)
	assert.Equal(mailSender.mail[0].address, "user@mail.com")
}

func TestDoMailJob(t *testing.T) {
	defer mailSender.Reset()
	dbReader.numUsers = 2
	assert := assert.New(t)
	err := MailJob(context.TODO(), dbReader, mailSender, now)
	assert.NoError(err)
	assert.Len(mailSender.mail, 2)
}
