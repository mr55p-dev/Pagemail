package mail

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/stretchr/testify/assert"
)

var (
	created  = time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	now      = time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)
	testUser = queries.ReadUsersWithMailRow{
		ID:       "123",
		Username: "user",
		Email:    "user@mail.com",
	}
	mailSender = &NoOpSender{}
)

func TestSendMailToUser(t *testing.T) {
	assert := assert.New(t)
	defer mailSender.Reset()
	db, mock, err := sqlmock.New()
	assert.NoError(err)
	mock.
		ExpectQuery(`.*`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"user_id",
				"url",
				"title",
				"description",
				"image_url",
				"preview_state",
				"created",
				"updated",
			}).
				AddRow(
					"aaa",
					testUser.ID,
					"https://example.com",
					nil,
					nil,
					nil,
					"unknown",
					time.Now(),
					time.Now(),
				),
		)

	err = SendUserDigest(context.TODO(), &testUser, db, mailSender, now)

	assert.NoError(err)
	assert.NoError(mock.ExpectationsWereMet())

	assert.Len(mailSender.mail, 1)
	assert.Equal(mailSender.mail[0].To, testUser.Email)
}

func TestDoMailJob(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	assert.NoError(err)

	mock.
		ExpectQuery(`.*`).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"id", "username", "email"}).
				AddRow(testUser.ID, testUser.Username, testUser.Email),
		)
	mock.
		ExpectQuery(`.*`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"user_id",
				"url",
				"title",
				"description",
				"image_url",
				"preview_state",
				"created",
				"updated",
			}).
				AddRow(
					"aaa",
					testUser.ID,
					"https://example.com",
					nil,
					nil,
					nil,
					"unknown",
					time.Now(),
					time.Now(),
				),
		)

	err = DigestJob(context.TODO(), db, mailSender, now)
	assert.NoError(err)
	assert.NoError(mock.ExpectationsWereMet())
	assert.Len(mailSender.mail, 1)
}
