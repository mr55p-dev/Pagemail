package mail

import  "context"

type TestClient struct{}

func (*TestClient) Send(ctx context.Context, user *User, body string) error {
	logger.Info("Test sending mail", "user", user.Email)
	return nil
}
