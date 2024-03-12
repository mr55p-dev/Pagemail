package mail

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/render"
)

func GenerateMailBody(ctx context.Context, user *User, pages []db.Page, since time.Time) ([]byte, error) {
	dest := new(bytes.Buffer)
	err := render.MailDigest(&since, user.Name, pages).Render(ctx, dest)
	return dest.Bytes(), err
}

type MailClient interface {
	SendMail(context.Context, *logging.Logger, *User, string) error
}

type SesMailClient struct {
	FromAddr  string
	SesClient *ses.Client
	log       *logging.Logger
}

func NewSesMailClient(ctx context.Context, log *logging.Logger) *SesMailClient {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}
	client := ses.NewFromConfig(cfg)
	return &SesMailClient{
		log:       log,
		SesClient: client,
		FromAddr:  "mail@pagemail.io",
	}
}

func (c *SesMailClient) SendMail(ctx context.Context, log *logging.Logger, user *User, body string) error {
	c.log.DebugContext(ctx, "Sending mail")
	_, err := c.SesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{user.Email},
		},
		Source: aws.String(c.FromAddr),
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    &body,
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &types.Content{
				Data:    aws.String("pagemail digest"),
				Charset: aws.String("UTF-8"),
			},
		},
		ReplyToAddresses: []string{c.FromAddr},
		Tags: []types.MessageTag{
			{
				Name:  aws.String("purpose"),
				Value: aws.String("daily-digest"),
			},
		},
	})
	if err != nil {
		c.log.Errc(ctx, "Error sending mail", err)
		return err
	}
	c.log.InfoContext(ctx, "Sent mail")
	return nil
}
