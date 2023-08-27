package readability

import (
	"log"
	"pagemail/server/models"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/spf13/cobra"
)

// Function which crawls all pages in the database
// Something like page preview hook
func CrawlAll(app *pocketbase.PocketBase, cfg models.ReaderConfig) *cobra.Command {
	cmd := func(c *cobra.Command, args []string) {
		ids := []models.Page{}
		err := app.Dao().DB().Select("id").From("pages").Build().All(&ids)
		if err != nil || len(ids) == 0 {
			log.Print(err)
			log.Print("No ids found")
			return
		}
		log.Printf("Collected %d ids", len(ids))

		for _, id := range ids {
			log.Print(id.Id)
			record, err := app.Dao().FindRecordById("pages", id.Id)
			last_crawled := record.GetDateTime("last_crawled")
			if !last_crawled.IsZero() {
				log.Printf("record %s has been crawled at %s", record.Id, last_crawled)
				continue
			}
			log.Printf("record %s has never been crawled (last crawled %s)", record.Id, last_crawled)

			url := record.GetString("url")
			res, err := FetchPreview(url, cfg)
			if err != nil {
				log.Panic(err)
			}
			form := forms.NewRecordUpsert(app, record)
			if err := form.LoadData(res.ToMap()); err != nil {
				log.Panic(err)
			}
			if err := form.Submit(); err != nil {
				log.Panic(err)
			}
			log.Printf("Finished updates for %s", id.Id)
		}
		return
	}

	out := cobra.Command{
		Use: "scrape-all",
		Run: cmd,
	}

	return &out
}
