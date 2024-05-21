package mail

import (
	"context"
	"io"
	"strings"
)

// NoOpSender implements MailSender for tests
type NoOpSender struct {
	mail []Message
}

// Send implements MailSender and stores messages in the instance
func (m *NoOpSender) Send(ctx context.Context, msg Message) error {
	dst := strings.Builder{}
	io.Copy(&dst, msg.Body)
	cnts := dst.String()
	logger.DebugCtx(ctx, "No-op mail send", "from", msg.From, "to", msg.To, "subject", msg.Subject, "contents", cnts)
	m.mail = append(m.mail, msg)
	return nil
}

// Reset clears the mail log for the mock instance
func (m *NoOpSender) Reset() {
	m.mail = make([]Message, 0)
}
