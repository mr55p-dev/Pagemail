package mail

import (
	"context"
)

type TestClient struct{}

func (*TestClient) SendMail(ctx context.Context, user *User, body string) error {
	log.Info("Test sending mail", "user", user.Email)
	return nil
}
