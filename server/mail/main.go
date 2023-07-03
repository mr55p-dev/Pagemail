package mail

import (
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type UrlRecord struct{ Url string }

func format_time(tm time.Time) string {
	return fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d-00:00",
		tm.Year(),
		tm.Month(),
		tm.Day(),
		tm.Hour(),
		tm.Minute(),
		tm.Second(),
	)
}

func get_mail_body(records []*models.Record, name string) string {
	var list_contents string
	for _, record := range records {
		url := record.GetString("url")
		created := record.GetTime("created")
		list_contents += fmt.Sprintf(`<li><a href="%s">%s</a> (created %s)</li>`, url, url, created)
	}

	res := fmt.Sprintf(`
<!DOCTYPE html>
<html>

<head>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/alvaromontoro/almond.css@latest/dist/almond.min.css" />
</head>

<body>
	<h1>Your saved pages</h1>
		<p>Hello, %s. Here are all the pages you have recently saved:</p>
	<ul>
		%s
	</ul>
</body>
</html>
	`, name, list_contents)
	return res
}

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
		log.Print(err)
		return err
	}

	log.Printf("Fetched %d users", len(users))

	// For each user
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour).Add(7 * time.Hour)
	now := time.Now()

	for _, usr := range users {
		//// Fetch all records which have created BETWEEN now-24hrs AND now
		// Write a query which selects all the records where user_id == usr.id and created > 7am yesterday
		records, err := app.Dao().FindRecordsByExpr("pages",
			dbx.HashExp{"user_id": usr.Id},
			dbx.NewExp("created BETWEEN {:start} AND {:end}", dbx.Params{"start": yesterday, "end": now}))
		if err != nil {
			log.Print(err)
			return err
		}
		if len(records) == 0 {
			log.Printf("Skipping %s", usr.Email)
			continue
		}
		log.Printf("Found %d records", len(records))

		//// Send an email with the links
		message := mailer.Message{
			From: mail.Address{
				Address: app.Settings().Meta.SenderAddress,
				Name:    app.Settings().Meta.SenderName,
			},
			To:      []mail.Address{{Address: usr.Email}},
			Subject: "Pagemail briefing",
			HTML:    get_mail_body(records, usr.Email),
		}
		log.Printf("Sending email to %s (%d links)", usr.Email, len(records))
		log.Print(get_mail_body(records, usr.Email))
		if err := mail_client.Send(&message); err != nil {
			log.Printf("Failed to send email to %s", usr.Email)
			log.Print(err)
		}
	}
	return nil
}
