package readability

import (
	"context"
	"fmt"
	"log"
	"pagemail/server/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

type JobCompletion struct {
	status *polly.GetSpeechSynthesisTaskOutput
	err    error
}

func FetchJobData(c context.Context, cfg *models.PMContext, id string) (*polly.GetSpeechSynthesisTaskOutput, error) {
	service := polly.NewFromConfig(*cfg.AWS)
	return service.GetSpeechSynthesisTask(c, &polly.GetSpeechSynthesisTaskInput{
		TaskId: &id,
	})
}

func AwaitJobCompletion(c context.Context, cfg *models.PMContext, id *string) (chan *polly.GetSpeechSynthesisTaskOutput, chan error) {
	taskData := make(chan *polly.GetSpeechSynthesisTaskOutput)
	taskErr := make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				taskErr <- fmt.Errorf("AwaitJobCompletion: PANIC with val %s", err)
			}
		}()

		for i := 0; i < 100; i++ {
			status, err := FetchJobData(context.TODO(), cfg, *id)
			if status == nil || err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			log.Printf(
				"Found task %s with status %s",
				*status.SynthesisTask.TaskId,
				string(status.SynthesisTask.TaskStatus),
			)

			switch status.SynthesisTask.TaskStatus {
			case types.TaskStatusCompleted:
				log.Printf("Status completed %s", status.SynthesisTask.TaskStatus)
				taskData <- status
				return
			case types.TaskStatusScheduled, types.TaskStatusInProgress:
				log.Printf("Status scheduled or in progress %s", status.SynthesisTask.TaskStatus)
				time.Sleep(time.Second * 2)
			case types.TaskStatusFailed:
				taskErr <- fmt.Errorf("Synthesis failed: %s", *status.SynthesisTask.TaskStatusReason)
				return
			default:
				log.Printf("Non-standard task status found: %s (reason %s)", status.SynthesisTask.TaskStatus, *status.SynthesisTask.TaskStatusReason)
			}
		}
		log.Print("Run out of retries, writing to task error chan")
		taskErr <- fmt.Errorf("Max number of retries reached polling AWS for task %s", *id)
	}()
	return taskData, taskErr
}


// // You can edit this code!
// // Click here and start typing.
// package main
//
// import (
// 	"fmt"
// 	"log"
// 	"time"
// )
//
// type Status string
//
// const (
// 	Complete   Status = "Complete"
// 	Failed     Status = "Failed"
// 	InProgress Status = "InProgress"
// 	Scheduled  Status = "Scheduled"
// )
//
// func main() {
//
// 	taskData := make(chan Status)
// 	taskErr := make(chan error)
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			status := Scheduled
//
// 			switch status {
// 			case Complete:
// 				log.Printf("Status completed %s", status)
// 				taskData <- status
// 				return
// 			case Scheduled, InProgress:
// 				log.Printf("Status scheduled or in progress %s", status)
// 				time.Sleep(time.Second * 2)
// 			case Failed:
// 				taskErr <- fmt.Errorf("Synthesis failed: %s", status)
// 				return
// 			default:
// 				log.Printf("Non-standard task status found: %s", status)
// 			}
// 		}
// 		log.Print("Run out of retries, writing to task error chan")
// 		taskErr <- fmt.Errorf("Max number of retries reached polling AWS for task ")
// 	}()
//
// 	select {
// 	case out := <-taskData:
// 		log.Print("Done", out)
// 	case <-taskErr:
// 		log.Print("Error")
// 	}
//
// }
