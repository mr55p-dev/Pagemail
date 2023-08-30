package readability

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
)

type JobCompletion struct {
	isSuccess bool
	status    *polly.GetSpeechSynthesisTaskOutput
	err       error
}

func FetchJobStatus(c context.Context, id string) (*polly.GetSpeechSynthesisTaskOutput, error) {
	cfg, err := config.LoadDefaultConfig(c)
	if err != nil {
		log.Panic(err)
	}

	service := polly.NewFromConfig(cfg)
	task, err := service.GetSpeechSynthesisTask(c, &polly.GetSpeechSynthesisTaskInput{
		TaskId: &id,
	})
	if err != nil {
		// Need to fall back to checking S3 here if the task is done
		log.Panic(err)
		return nil, err
	}
	return task, nil
}

func AwaitJobCompletion(c context.Context, id *string) chan JobCompletion {
	out := make(chan JobCompletion)
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for i := 0; i < 10; i++ {
			status, err := FetchJobStatus(context.TODO(), *id)
			if status == nil && err != nil {
				log.Printf("Cancelling inner context, status: %s, err: %s", string(status.SynthesisTask.TaskStatus), err)
				time.Sleep(time.Second)
			} else {
				out <- JobCompletion{
					isSuccess: true,
					status:    status,
					err:       nil,
				}
				return
			}
		}

		out <- JobCompletion{
			isSuccess: false,
			status:    nil,
			err:       fmt.Errorf("Max number of retries reached polling AWS for task %s", *id),
		}
	}()
	return out
}
