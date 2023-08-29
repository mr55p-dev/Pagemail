package readability

import (
	"log"
	"pagemail/server/models"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
)



func UpdateJobState(app *pocketbase.PocketBase, pageId string, state models.ReadabilityStatus) error {
	log.Printf("Looking for page wih id %s", pageId)
	record, err := app.Dao().FindRecordById("pages", pageId)
	if err != nil {
		log.Print(err)
		return err
	}

	form := forms.NewRecordUpsert(app, record)
	if err := form.LoadData(map[string]any{
		"readability_status": state,
	}); err != nil {
		return err
	}
	return form.Submit()
}
