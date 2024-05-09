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

// SesMailClient implements the MailSender interface using AWS SES SDK
type SesMailClient struct {
	sesClient *ses.Client
}

func NewSesMailClient(ctx context.Context, awsConfig aws.Config) *SesMailClient {
	client := ses.NewFromConfig(awsConfig)
	return &SesMailClient{
		sesClient: client,
	}
}

// Send will produce an email to the given address with the given body
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
		Source: aws.String(FROM_ADDR),
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
		ReplyToAddresses: []string{FROM_ADDR},
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
