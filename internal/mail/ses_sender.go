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

// AwsSender implements the MailSender interface using AWS SES SDK
type AwsSender struct {
	sesClient *ses.Client
}

func NewAwsSender(ctx context.Context, awsConfig aws.Config) *AwsSender {
	client := ses.NewFromConfig(awsConfig)
	return &AwsSender{
		sesClient: client,
	}
}

// Send will produce an email to the given address with the given body
func (c *AwsSender) Send(ctx context.Context, msg Message) error {
	logger := logger.With("recipient", msg.Body)
	logger.DebugCtx(ctx, "Sending mail", "recipient")
	bodyTextBuilder := strings.Builder{}
	_, err := io.Copy(&bodyTextBuilder, msg.Body)
	if err != nil {
		return fmt.Errorf("Failed to copy mail body: %w", err)
	}

	tags := make([]types.MessageTag, 0)
	for _, v := range msg.Tags {
		tags = append(tags, types.MessageTag{
			Name:  &v.Name,
			Value: &v.Value,
		})
	}

	bodyText := bodyTextBuilder.String()
	_, err = c.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{msg.To},
		},
		Source: &msg.From,
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    &bodyText,
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &types.Content{
				Data:    &msg.Subject,
				Charset: aws.String("UTF-8"),
			},
		},
		ReplyToAddresses: []string{ADDR_DIGEST},
		Tags:             tags,
	})
	if err != nil {
		logger.WithError(err).ErrorCtx(ctx, "Error sending mail")
		return err
	}
	logger.DebugCtx(ctx, "Sent mail")
	return nil
}
