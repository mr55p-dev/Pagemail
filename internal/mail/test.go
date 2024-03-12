package mail

import (
	"context"

	"github.com/mr55p-dev/pagemail/internal/logging"
)

type TestClient struct{}

func (*TestClient) SendMail(ctx context.Context, log *logging.Logger, user *User, body string) error {
	log.Info("Test sending mail", "user", user.Email)
	return nil
}
