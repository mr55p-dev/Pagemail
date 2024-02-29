package mail

import (
	"context"
	"os"

	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
)

func mapDbUser(u *db.User) User {
	return User{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}

func sendTestEmail() {
	log := logging.NewVoid()
	dbClient := db.NewClient("db/test.sqlite3", log)
	mailClient := NewSesMailClient(context.Background(), log)
	user, _ := dbClient.ReadUserById(context.Background(), "iy17XbjTy7")
	log.Info("Got user", "user", user)
	u := mapDbUser(user)
	pages, _ := dbClient.ReadPagesByUserId(context.Background(), user.Id, 1)
	since := Yesterday()
	msg, err := GenerateMailBody(context.Background(), &u, pages, since)
	os.WriteFile("html", msg, 0o777)
	err = mailClient.SendMail(context.Background(), log, &u, string(msg))
	if err != nil {
		log.Err("Error sending", err)
	} else {
		log.Info("Sent mail")
	}

}

func testUserGeneration() {
	log := logging.NewVoid()
	dbClient := db.NewClient("db/test.sqlite3", log)
	mailClient := new(TestClient)
	err := DoDigestJob(context.Background(), log, dbClient, mailClient)
	if err != nil {
		log.Err("Error sending", err)
	} else {
		log.Info("Done")
	}
}
