package mail

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SesMailClient struct {
	sesClient *ses.Client
}

func NewSesMailClient(ctx context.Context, awsConfig aws.Config) *SesMailClient {
	client := ses.NewFromConfig(awsConfig)
	return &SesMailClient{
		sesClient: client,
	}
}

func (c *SesMailClient) Send(ctx context.Context, recipient, body string) error {
	logger := logger.With("recipient", body)
	logger.DebugCtx(ctx, "Sending mail", "recipient")
	_, err := c.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Source: aws.String(from_addr),
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    &body,
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &types.Content{
				Data:    aws.String("Pagemail digest"),
				Charset: aws.String("UTF-8"),
			},
		},
		ReplyToAddresses: []string{from_addr},
		Tags: []types.MessageTag{
			{
				Name:  aws.String("purpose"),
				Value: aws.String("daily-digest"),
			},
		},
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Error sending mail")
		return err
	}
	logger.DebugCtx(ctx, "Sent mail")
	return nil
}
