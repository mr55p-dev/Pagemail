package mail

import (
	"context"
	"time"

	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

// MailDbReaderNoOp implements the MailDbReader interface with mock data
type MailDbReaderNoOp struct {
	created  time.Time
	user     dbqueries.User
	numPages int
	numUsers int
}

// ReadUsersWithMail passes mock data to the SES mailer function
func (m *MailDbReaderNoOp) ReadUsersWithMail(context.Context) ([]dbqueries.User, error) {
	users := make([]dbqueries.User, 0)
	for i := 0; i < m.numUsers; i++ {
		users = append(users, testUser)
	}
	return users, nil
}

// ReadPagesByUserBetween passes mock data for a user
func (m *MailDbReaderNoOp) ReadPagesByUserBetween(context.Context, dbqueries.ReadPagesByUserBetweenParams) ([]dbqueries.Page, error) {
	testPages := []dbqueries.Page{
		{
			ID:      "123",
			UserID:  m.user.ID,
			Url:     "https://pagemail.io",
			Created: m.created,
		},
		{
			ID:      "123",
			UserID:  m.user.ID,
			Url:     "https://example.com",
			Created: m.created,
		},
	}
	if m.numPages == 0 {
		return []dbqueries.Page{}, nil
	}
	return testPages[:m.numPages], nil
}
