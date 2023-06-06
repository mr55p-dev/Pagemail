package custom_api

import (
	"fmt"
	"log"
	"net/mail"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type User struct {
	Id    string
	Email string
}

func Mailer(app *pocketbase.PocketBase) error {
	mail_client := app.NewMailClient()
	db := app.Dao().DB()

	log.Print("Running mailer")

	// Fetch all the subscribed users
	var users []User
	q := db.NewQuery("SELECT id, email FROM users WHERE subscribed = true AND email IS NOT NULL")
	err := q.All(&users)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Fetched %d users", len(users))

	// For each user
	for _, usr := range users {
		log.Printf("User ID: %s", usr.Id)
		log.Print(usr)
		//// Fetch all records which have created BETWEEN now-24hrs AND now
		var urls []struct{Url string}
		q_str := fmt.Sprintf("SELECT url FROM pages WHERE user_id = '%s'", usr.Id)
		q := db.NewQuery(q_str)
		err := q.All(&urls)
		if err != nil {
			log.Fatal(err)
			return err
		}
		if len(urls) == 0 {
			continue
		}
		//// ForEach record => Preview

		//// Send an email with the links
		message := mailer.Message{
			From: mail.Address{
				Address: app.Settings().Meta.SenderAddress,
				Name:    app.Settings().Meta.SenderName,
			},
			To:      []mail.Address{{Address: usr.Email}},
			Subject: "Pagemail briefing",
			HTML:    "<h1>This is a test message</h1>",
		}
		log.Printf("Sending test email to %s", message.To[0])
		if err := mail_client.Send(&message); err != nil {
			log.Print("Failed to send email")
			log.Print(err)
		}

	}
	return nil
}
