package mail

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/mr55p-dev/pagemail/internal/db"
)

// Message holds basic properties about generated emails
type Message struct {
	address  string
	contents string
}

// MailSenderNoOp implements MailSender for tests
type MailSenderNoOp struct {
	mail []Message
}

// Send implements MailSender and stores messages in the instance
func (m *MailSenderNoOp) Send(ctx context.Context, addr string, contents io.Reader) error {
	dst := strings.Builder{}
	io.Copy(&dst, contents)
	cnts := dst.String()
	logger.DebugCtx(ctx, "No-op mail send", "address", addr, "contents", cnts)
	m.mail = append(m.mail, Message{
		address:  addr,
		contents: cnts,
	})
	return nil
}

// Reset clears the mail log for the mock instance
func (m *MailSenderNoOp) Reset() {
	m.mail = make([]Message, 0)
}

// MailDbReaderNoOp implements the MailDbReader interface with mock data
type MailDbReaderNoOp struct {
	created  time.Time
	user     db.User
	numPages int
	numUsers int
}

// ReadUsersWithMail passes mock data to the SES mailer function
func (m *MailDbReaderNoOp) ReadUsersWithMail(context.Context) ([]db.User, error) {
	users := make([]db.User, 0)
	for i := 0; i < m.numUsers; i++ {
		users = append(users, testUser)
	}
	return users, nil
}

// ReadPagesByUserBetween passes mock data for a user
func (m *MailDbReaderNoOp) ReadPagesByUserBetween(context.Context, string, time.Time, time.Time) ([]db.Page, error) {
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
