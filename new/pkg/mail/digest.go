package mail

import (
	"context"

	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

func DoDigestJob(ctx context.Context, dbClient *db.Client, mailClient *MailClient) error {
	users, err := dbClient.ReadUsersWithMail(ctx)
	if err != nil {
		log.ErrContext(ctx, "Failed digest job", err)
		return err
	}
	for _, user := range users {
		log.InfoContext(ctx, "Looking up pages", logging.UserId, user.Id)
		pages, err := dbClient.ReadPagesByUserId(ctx, user.Id)
		if err != nil {
			log.ErrContext(ctx, "Failed looking up pages", err)
		}
		message, err := GenerateMailBody(ctx, &user, pages)
		if err != nil {
			log.ErrContext(ctx, "Failed while generating mail body", err)
			return err
		}
		err = mailClient.SendMail(ctx, user, string(message))
		if err != nil {
			log.ErrContext(ctx, "Failed sending mail", err)
			return err
		}
	}
	log.Info("Mail send job done")
	return nil
}
