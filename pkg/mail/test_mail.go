package mail

import (
	"context"
	"os"

	"github.com/mr55p-dev/pagemail/pkg/db"
)

func sendTestEmail() {
	dbClient := db.NewClient("db/test.sqlite3")
	mailClient := NewSesMailClient(context.Background())
	user, _ := dbClient.ReadUserById(context.Background(), "iy17XbjTy7")
	log.Info("Got user", "user", user)
	pages, _ := dbClient.ReadPagesByUserId(context.Background(), user.Id, -1)
	msg, err := GenerateMailBody(context.Background(), user, pages)
	os.WriteFile("html", msg, 0o777)
	err = mailClient.SendMail(context.Background(), user, string(msg))
	if err != nil {
		log.Err("Error sending", err)
	} else {
		log.Info("Sent mail")
	}

}

func testUserGeneration() {
	dbClient := db.NewClient("db/test.sqlite3")
	mailClient := new(TestClient)
	err := DoDigestJob(context.Background(), dbClient, mailClient)
	if err != nil {
		log.Err("Error sending", err)
	} else {
		log.Info("Done")
	}
}
