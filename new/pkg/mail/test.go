package mail

import (
	"context"

	"github.com/mr55p-dev/pagemail/pkg/db"
)

type TestClient struct{}

func (*TestClient) SendMail(ctx context.Context, user *db.User, body string) error {
	log.Info("Test sending mail", "user", user.Email, "body", body)
	return nil
}
