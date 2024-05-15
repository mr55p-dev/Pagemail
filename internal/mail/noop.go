package mail

import (
	"context"
	"io"
	"strings"
	"time"
)

// msg holds basic properties about generated emails
type msg struct {
	address  string
	contents string
}

// NoOpSender implements MailSender for tests
type NoOpSender struct {
	mail []msg
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
func (m *NoOpSender) Send(ctx context.Context, addr string, contents io.Reader) error {
	dst := strings.Builder{}
	io.Copy(&dst, contents)
	cnts := dst.String()
	logger.DebugCtx(ctx, "No-op mail send", "address", addr, "contents", cnts)
	m.mail = append(m.mail, msg{
		address:  addr,
		contents: cnts,
	})
	return nil
}

// Reset clears the mail log for the mock instance
func (m *NoOpSender) Reset() {
	m.mail = make([]msg, 0)
}
