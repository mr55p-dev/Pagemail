package mail

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type UrlRecord struct{ Url string }

type MailItem struct {
	Url       string
	Timestamp time.Time
}

type User struct {
	Id    string
	Email string
}

type MailTemplateData struct {
	Name string
	SavedItems []MailItem
}

func get_mail_body(w io.Writer, items []MailItem, name string) string {
	tmpl := template.Must(template.ParseFiles("templates/mailer.html.template"))
	tmpl.Execute(w, MailT)
	var list_contents string
	for _, item := range items {
		url := item.Url
		created := item.Timestamp.Format("02-01 15:04")
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

func TestMailer(c echo.Context) error {
	mailItems := []MailItem{
		{
			Url:       "http://www.example.com/",
			Timestamp: time.Now(),
		},
		{
			Url:       "https://mail.google.com/this_is_a_long_url/with?some=query&parameters=to&be=annoying",
			Timestamp: time.Now(),
		},
		{
			Url:       "http://pagemail.io/",
			Timestamp: time.Now(),
		},
	}
	mailHTML := get_mail_body(mailItems, "Test user")
	return c.HTML(http.StatusOK, mailHTML)
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

		mailItems := []MailItem{}
		for _, record := range records {
			mailItems = append(mailItems, MailItem{
				Url:       record.GetString("url"),
				Timestamp: record.GetTime("created"),
			})
		}

		//// Send an email with the links
		message := mailer.Message{
			From: mail.Address{
				Address: app.Settings().Meta.SenderAddress,
				Name:    app.Settings().Meta.SenderName,
			},
			To:      []mail.Address{{Address: usr.Email}},
			Subject: "Pagemail briefing",
			HTML:    get_mail_body(mailItems, usr.Email),
		}
		log.Printf("Sending email to %s (%d links)", usr.Email, len(records))

		if err := mail_client.Send(&message); err != nil {
			log.Printf("Failed to send email to %s", usr.Email)
			log.Print(err)
		}
	}
	return nil
}
