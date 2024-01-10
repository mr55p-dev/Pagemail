package mail

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
	"github.com/mr55p-dev/pagemail/pkg/render"
)

var log logging.Log

const (
	Sender = "ellis@pagemail.io"
)

func init() {
	log = logging.GetLogger("mail")
}

func GenerateMailBody(ctx context.Context, user *db.User, pages []db.Page) ([]byte, error) {
	now := time.Now()
	dt := time.Date(now.Year(), now.Month(), now.Day()-1, 7, 0, 0, 0, time.UTC)
	dest := new(bytes.Buffer)
	err := render.MailDigest(&dt, user, pages).Render(ctx, dest)
	return dest.Bytes(), err
}

type MailClient interface {
	SendMail(context.Context, *db.User, string) error
}

type SesMailClient struct {
	FromAddr  string
	SesClient *ses.Client
}

func NewSesMailClient(ctx context.Context) *SesMailClient {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}
	client := ses.NewFromConfig(cfg)
	return &SesMailClient{
		SesClient: client,
	}
}

func (c *SesMailClient) SendMail(ctx context.Context, user *db.User, body string) error {
	log.Debug("Sending mail", logging.UserMail, user.Email)
	out, err := c.SesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{user.Email},
		},
		Source: &c.FromAddr,
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
		ReplyToAddresses: []string{Sender},
		Tags: []types.MessageTag{
			{
				Name:  aws.String("purpose"),
				Value: aws.String("daily-digest"),
			},
		},
	})
	if err != nil {
		log.ErrorContext(ctx, "Error sending mail", logging.Error, err.Error())
		return err
	}
	log.InfoContext(ctx, "Sent mail", logging.UserMail, user.Email, "message-id", out.MessageId)
	return nil
}
