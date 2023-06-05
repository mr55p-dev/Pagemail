package custom_api

import (
	"log"
	"net/mail"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

func Mailer(app *pocketbase.PocketBase) {
	mail_client := app.NewMailClient()
	message := mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To:      []mail.Address{{Address: "ellislunnon@gmail.com"}},
		Subject: "TEST_SUBJECT",
		HTML:    "<h1>This is a test message</h1>",
	}
	log.Printf("Sending test email to %s", message.To[0])
	if err := mail_client.Send(&message); err != nil {
		log.Print("Failed to send email")
		log.Print(err)
	}
}
