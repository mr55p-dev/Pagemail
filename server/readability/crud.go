package readability

import (
	"log"
	"encoding/json"
	"pagemail/server/models"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
)



func UpdateJobState(app *pocketbase.PocketBase, pageId string, state models.ReadabilityStatus, taskData *polly.StartSpeechSynthesisTaskOutput) error {
	record, err := app.Dao().FindRecordById("pages", pageId)
	if err != nil {
		log.Print(err)
		return err
	}

	form := forms.NewRecordUpsert(app, record)
	newData := map[string]any{
		"readability_status": state,
	}
	if taskData != nil {
		data, err := json.Marshal(taskData)
		if err != nil {
			return err
		}
		newData["readability_task_data"] = string(data)
	}
	if err := form.LoadData(newData); err != nil {
		return err
	}
	return form.Submit()
}

// func GetPage(app *pocketbase.PocketBase, pageId string) error {
// 	record, err := app.Dao().FindRecordById("pages", pageId)
// 	if err != nil {
// 		return err
// 	}
// }
