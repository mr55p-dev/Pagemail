package mail

import (
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/pocketbase/pocketbase"
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

func get_mail_body(urls *[]UrlRecord, name string) string {
	var list_contents string
	for _, url := range *urls {
		list_contents += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, url.Url, url.Url)
	}

	res := fmt.Sprintf(`
<!DOCTYPE html>
<html>

<head>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/alvaromontoro/almond.css@latest/dist/almond.min.css" />
</head>

<body>
	<h1>Your Pagemail digest</h1>
	<p>Hello, %s. Here are all the pages you have recently saved</p>
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
	day_duration, _ := time.ParseDuration("24h")
	now := time.Now()
	start := now.Add(day_duration * -1)

	for _, usr := range users {
		log.Printf("User ID: %s", usr.Id)
		log.Print(usr)
		//// Fetch all records which have created BETWEEN now-24hrs AND now
		var urls []UrlRecord
		q_str := fmt.Sprintf("SELECT url FROM pages WHERE user_id = '%s' AND created BETWEEN '%s' AND '%s'", usr.Id, format_time(start), format_time(now))
		q := db.NewQuery(q_str)
		err := q.All(&urls)
		if err != nil {
			log.Print(err)
			return err
		}
		if len(urls) == 0 {
			continue
		}
		log.Print(urls)

		//// Send an email with the links
		message := mailer.Message{
			From: mail.Address{
				Address: app.Settings().Meta.SenderAddress,
				Name:    app.Settings().Meta.SenderName,
			},
			To:      []mail.Address{{Address: usr.Email}},
			Subject: "Pagemail briefing",
			HTML:    get_mail_body(&urls, usr.Email),
		}
		log.Printf("Sending test email to %s", message.To[0])
		if err := mail_client.Send(&message); err != nil {
			log.Printf("Failed to send email to %s", usr.Email)
			log.Print(err)
		}
	}
	return nil
}
