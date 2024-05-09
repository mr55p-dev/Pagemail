package mail

import (
	"context"
	"fmt"
	"io"
	"strings"

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

func (c *SesMailClient) Send(ctx context.Context, recipient string, body io.Reader) error {
	logger := logger.With("recipient", body)
	logger.DebugCtx(ctx, "Sending mail", "recipient")
	bodyTextBuilder := strings.Builder{}
	_, err := io.Copy(&bodyTextBuilder, body)
	if err != nil {
		return fmt.Errorf("Failed to copy mail body: %w", err)
	}

	bodyText := bodyTextBuilder.String()
	_, err = c.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Source: aws.String(from_addr),
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    &bodyText,
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
