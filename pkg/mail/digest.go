package mail

import (
	"context"

	"github.com/mr55p-dev/pagemail/pkg/db"
	"github.com/mr55p-dev/pagemail/pkg/logging"
)

func DoDigestJob(ctx context.Context, dbClient *db.Client, mailClient MailClient) error {
	log.InfoContext(ctx, "Starting mail send job")
	users, err := dbClient.ReadUsersWithMail(ctx)
	if err != nil {
		log.ErrContext(ctx, "Failed digest job", err)
		return err
	}
	for _, user := range users {
		log.DebugContext(ctx, "Looking up pages", logging.UserId, user.Id)
		pages, err := dbClient.ReadPagesByUserId(ctx, user.Id)
		if err != nil {
			log.ErrContext(ctx, "Failed looking up pages", err)
		}
		log.DebugContext(ctx, "Found pages", "count", len(pages))
		message, err := GenerateMailBody(ctx, &user, pages)
		if err != nil {
			log.ErrContext(ctx, "Failed while generating mail body", err)
			return err
		}
		log.DebugContext(ctx, "Generated message", "message", string(message))
		err = mailClient.SendMail(ctx, &user, string(message))
		if err != nil {
			log.ErrContext(ctx, "Failed sending mail", err)
			return err
		}
	}
	log.InfoContext(ctx, "Mail send job done")
	return nil
}
