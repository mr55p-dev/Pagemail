package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"pagemail/server/models"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type MailTemplateData struct {
	UserIdentifier string
	DateStart      string
	Pages          []models.UrlData
}

func GetUsers(db *dbx.Builder) ([]models.User, error) {
	var users []models.User
	q := (*db).NewQuery("SELECT id, email FROM users WHERE subscribed = true AND email IS NOT NULL")
	err := q.All(&users)
	return users, err
}

func GetUserPages(app *pocketbase.PocketBase, user models.User, startTime time.Time) ([]models.PageRecord, error) {
	//// Fetch all records which have created BETWEEN now-24hrs AND now
	// Write a query which selects all the records where user_id == usr.id and created > 7am yesterday
	var pages []models.PageRecord
	records, err := app.Dao().FindRecordsByExpr("pages",
		dbx.HashExp{"user_id": user.Id},
		dbx.NewExp("created BETWEEN {:start} AND {:end}", dbx.Params{"start": startTime, "end": time.Now()}))
	if err != nil {
		log.Print(err)
		return pages, err
	}
	for _, row := range records {
		pages = append(pages, models.PageRecord{
			Url:     row.GetString("url"),
			Created: row.GetCreated().Time(),
		})
	}

	return pages, nil
}

func GetPageData(page models.PageRecord) (models.UrlData, error) {
	res := models.UrlData{
		PageRecord: models.PageRecord{
			Url: page.Url,
		},
	}
	return res, nil
}

func GetMailBody(data MailTemplateData) string {
	templatePath := os.Getenv("PAGEMAIL_EMAIL_TEMPLATE_PATH")
	var w bytes.Buffer
	tmpl := template.Must(template.ParseFiles(templatePath))
	if err := tmpl.Execute(&w, data); err != nil {
		panic(err)
	}

	return w.String()
}

func getUserIdentifier(user models.User) string {
	if user.Name != "" {
		return user.Name
	} else if user.Email != "" {
		return user.Email
	} else {
		return user.Id
	}
}

func Mailer(app *pocketbase.PocketBase) error {
	log.Print("Running mailer")

	// Setup clients
	mail_client := app.NewMailClient()
	db := app.Dao().DB()

	// Fetch all the subscribed users
	users, err := GetUsers(&db)
	log.Printf("Fetched %d users", len(users))
	if err != nil {
		log.Print(fmt.Errorf("Failed to fetch users: %s", err))
		return err
	}

	// For each user
	startDate := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour).Add(7 * time.Hour)

	for _, usr := range users {
		// Get the relevant pages for the user
		pages, err := GetUserPages(app, usr, startDate)
		if err != nil {
			log.Print(fmt.Errorf("Failed to fetch pages for user %s: %s", usr.Id, err))
		}
		log.Printf("Found %d records", len(pages))

		// Enrich page data with previews
		var enrichedPages []models.UrlData
		for _, page := range pages {
			if d, err := GetPageData(page); err == nil {
				enrichedPages = append(enrichedPages, d)
			}
		}

		// Skip if the user does not have any pages after enriching
		if len(enrichedPages) == 0 {
			continue
		}

		// Create mail template data
		identifier := getUserIdentifier(usr)
		data := MailTemplateData{
			UserIdentifier: identifier,
			DateStart:      startDate.Format("02/01/06"),
			Pages:          enrichedPages,
		}

		// Send an email with the links
		message := mailer.Message{
			From: mail.Address{
				Address: app.Settings().Meta.SenderAddress,
				Name:    app.Settings().Meta.SenderName,
			},
			To:      []mail.Address{{Address: usr.Email}},
			Subject: "Pagemail briefing",
			HTML:    GetMailBody(data),
		}
		log.Printf("Sending email to %s (%d pages)", usr.Email, len(data.Pages))

		if err := mail_client.Send(&message); err != nil {
			log.Printf("Failed to send email to %s", usr.Email)
			log.Print(err)
		}
	}
	return nil
}

func TestMailBody(c echo.Context) error {
	templateData := MailTemplateData{
		UserIdentifier: "Test user",
		DateStart:      time.Now().Format("02-01-2006"),
		Pages: []models.UrlData{{
			UrlPreviewData: models.UrlPreviewData{
				Title:       "Example page",
				Description: "A disturbed garage's creature comes with it the thought that the urdy carnation is a flock. A cordial cushion's apartment comes with it the thought that the starboard umbrella is a weapon. Authors often misinterpret the chocolate as an advised college, when in actuality it feels more like a saintly dinghy.",
			},
			PageRecord: models.PageRecord{
				Url:     "https://www.example.com/",
				Created: time.Now(),
			},
		},
			{
				UrlPreviewData: models.UrlPreviewData{
					Title:       "Example page",
					Description: "",
				},
				PageRecord: models.PageRecord{
					Url:     "https://www.example.com/",
					Created: time.Now(),
				},
			},
			{
				PageRecord: models.PageRecord{
					Url:     "https://www.example2.com",
					Created: time.Now(),
				},
			},
		},
	}

	mailHTML := GetMailBody(templateData)
	return c.HTML(http.StatusOK, mailHTML)
}
