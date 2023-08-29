package readability

import (
	"context"
	"log"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
)

func FetchJobStatus(ctx context.Context, id string) (*string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Panic(err)
	}

	service := polly.NewFromConfig(cfg)
	task, err := service.GetSpeechSynthesisTask(ctx, &polly.GetSpeechSynthesisTaskInput{
		TaskId: &id,
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Print(string(task.SynthesisTask.TaskStatus))
	return task.SynthesisTask.TaskId, nil
}
