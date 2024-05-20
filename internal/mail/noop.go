package mail

import (
	"context"
	"io"
	"strings"
	"time"
)

// NoOpSender implements MailSender for tests
type NoOpSender struct {
	mail []Message
}

func NewNoOpSender(ctx context.Context) *NoOpSender {
	sender := new(NoOpSender)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Hour):
				sender.Reset()
			}
		}
	}()
	return sender
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
