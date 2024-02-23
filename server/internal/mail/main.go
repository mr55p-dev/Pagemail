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
	"github.com/mr55p-dev/pagemail/internal/render"
)

var log logging.Log

func init() {
	log = logging.GetLogger("mail")
}

func GenerateMailBody(ctx context.Context, user *User, pages []db.Page, since time.Time) ([]byte, error) {
	dest := new(bytes.Buffer)
	err := render.MailDigest(&since, user.Name, pages).Render(ctx, dest)
	return dest.Bytes(), err
}

type MailClient interface {
	SendMail(context.Context, *User, string) error
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
		FromAddr:  "mail@pagemail.io",
	}
}

func (c *SesMailClient) SendMail(ctx context.Context, user *User, body string) error {
	log.Debug("Sending mail", logging.UserMail, user.Email, "from", c.FromAddr)
	out, err := c.SesClient.SendEmail(ctx, &ses.SendEmailInput{
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
		log.ErrorContext(ctx, "Error sending mail", logging.Error, err.Error())
		return err
	}
	log.InfoContext(ctx, "Sent mail", logging.UserMail, user.Email, "message-id", out.MessageId)
	return nil
}
