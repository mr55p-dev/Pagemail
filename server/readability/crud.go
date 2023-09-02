package readability

import (
	"log"
	"encoding/json"
	"pagemail/server/models"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
)



func UpdateJobState(app *pocketbase.PocketBase, pageId string, state models.ReadabilityStatus, taskData *polly.StartSpeechSynthesisTaskOutput) {
	record, err := app.Dao().FindRecordById("pages", pageId)
	if err != nil {
		log.Panic(err)
	}

	form := forms.NewRecordUpsert(app, record)
	newData := map[string]any{
		"readability_status": state,
	}
	if taskData != nil {
		data, err := json.Marshal(taskData)
		if err != nil {
			log.Panic(err)
		}
		newData["readability_task_data"] = string(data)
		newData["readability_task_id"] = *taskData.SynthesisTask.TaskId
	}
	if err := form.LoadData(newData); err != nil {
		log.Panic(err)
	}
	err = form.Submit()
	if err != nil {
		log.Panic(err)
	}
	return
}

// func GetPage(app *pocketbase.PocketBase, pageId string) error {
// 	record, err := app.Dao().FindRecordById("pages", pageId)
// 	if err != nil {
// 		return err
// 	}
// }
