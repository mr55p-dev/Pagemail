package main

import (
	"context"
	"os"

	"github.com/mr55p-dev/pagemail/pkg/aws"
	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
	"github.com/mr55p-dev/pagemail/pkg/mail"
)

var log logging.Log

func sendTestEmail() {
	dbClient := db.NewClient()
	mailClient := &mail.SesMailClient{
		SesClient: aws.GetSesClient(context.Background()),
		FromAddr:  "ellis@pagemail.io",
	}
	user, _ := dbClient.ReadUserById(context.Background(), "iy17XbjTy7")
	log.Info("Got user", "user", user)
	pages, _ := dbClient.ReadPagesByUserId(context.Background(), user.Id)
	msg, err := mail.GenerateMailBody(context.Background(), user, pages)
	os.WriteFile("mail.html", msg, 0o777)
	err = mailClient.SendMail(context.Background(), user, string(msg))
	if err != nil {
		log.Err("Error sending", err)
	} else {
		log.Info("Sent mail")
	}

}

type TestClient struct{}

func (*TestClient) SendMail(ctx context.Context, user *db.User, body string) error {
	log.Info("Test sending mail", "user", user.Email, "body", body)
	return nil
}

func testUserGeneration() {
	dbClient := db.NewClient()
	mailClient := new(TestClient)
	err := mail.DoDigestJob(context.Background(), dbClient, mailClient)
	if err != nil {
		log.Err("Error sending", err)
	} else {
		log.Info("Done")
	}
}

func main() {
	log = logging.GetLogger("test")
	sendTestEmail()
	// testUserGeneration()
}
