package readability

import (
	"log"
	"pagemail/server/preview"
	"pagemail/server/readability"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	// "github.com/pocketbase/pocketbase/forms"
)

// Function which crawls all pages in the database
// Something like page preview hook
func CrawlAll(app *pocketbase.PocketBase, cfg *readability.ReaderConfig) error {
	// Select * from pages
	// go fetch preview (FetchPreview)
	// update record
	records, err := app.Dao().FindRecordsByExpr("pages", dbx.NewExp(""))
	if err != nil {
		return err
	}
	log.Printf("Collected %d records", len(records))

	// for _, record := range records {
	// 	url := record.Get("url")
	// 	if url == nil {
	// 		continue
	// 	}
	// 	res, err := preview.FetchPreview(url, cfg &readability.ReaderConfig)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	form := forms.NewRecordUpsert(app, record)
	// 	if err := form.LoadData(res.ToMap()); err != nil {
	// 		return err
	// 	}
	// 	if err := form.Submit(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
