package mail

import (
	"context"
	"os"

	"github.com/mr55p-dev/pagemail/internal/db"
)

func mapDbUser(u *db.User) User {
	return User{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}

func sendTestEmail() {
	dbClient := db.NewClient("db/test.sqlite3", nil)
	mailClient := NewSesMailClient(context.Background())
	user, _ := dbClient.ReadUserById(context.Background(), "iy17XbjTy7")
	logger.Info("Got user", "user", user)
	u := mapDbUser(user)
	pages, _ := dbClient.ReadPagesByUserId(context.Background(), user.Id, 1)
	since := Yesterday()
	msg, err := GenerateMailBody(context.Background(), &u, pages, since)
	os.WriteFile("html", msg, 0o777)
	err = mailClient.Send(context.Background(), &u, string(msg))
	if err != nil {
		logger.Error("Error sending", err)
	} else {
		logger.Info("Sent mail")
	}

}

func testUserGeneration() {
	dbClient := db.NewClient("db/test.sqlite3", nil)
	mailClient := new(TestClient)
	err := DoDigestJob(context.Background(), dbClient, mailClient)
	if err != nil {
		logger.Error("Error sending", err)
	} else {
		logger.Info("Done")
	}
}
