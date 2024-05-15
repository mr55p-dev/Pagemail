package mail

import (
	"context"
	"testing"
	"time"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/stretchr/testify/assert"
)


var (
	created  = time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	now      = time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)
	testUser = dbqueries.User{
		ID:         "123",
		Username:   "user",
		Email:      "user@mail.com",
		Subscribed: true,
	}

	mailSender = &NoOpSender{}
	dbReader   = &MailDbReaderNoOp{
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
