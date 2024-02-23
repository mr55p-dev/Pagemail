package mail

import (
	"context"
	"sync"
	"time"

	"github.com/mr55p-dev/pagemail/internal/db"
)

type User struct {
	Id    string
	Name  string
	Email string
}

func GetUsers(ctx context.Context, dbClient *db.Client) (out []User, err error) {
	users, err := dbClient.ReadUsersWithMail(ctx)
	if err != nil {
		return
	}
	for _, v := range users {
		out = append(out, User{
			Id:    v.Id,
			Name:  v.Name,
			Email: v.Email,
		})
	}
	return
}

func Yesterday() time.Time {
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day()-1, 7, 0, 0, 0, time.UTC)
	return since
}

func GetEmailForUser(ctx context.Context, dbClient *db.Client, user User) (message string, err error) {
	log.DebugContext(ctx, "Looking up pages", logging.UserId, user)
	since := Yesterday()

	pages, err := dbClient.ReadPagesByUserId(ctx, user.Id, -1)
	if err != nil {
		log.ErrContext(ctx, "Failed looking up pages", err)
	}
	pages = FilterPages(ctx, pages, since)
	log.DebugContext(ctx, "Found pages", "count", len(pages))
	messageBytes, err := GenerateMailBody(ctx, &user, pages, since)
	if err != nil {
		log.ErrContext(ctx, "Failed while generating mail body", err)
		return
	}
	message = string(messageBytes)
	return
}

func FilterPages(ctx context.Context, pages []db.Page, since time.Time) (out []db.Page) {
	for _, v := range pages {
		if v.Created.After(since) {
			out = append(out, v)
		}
	}
	return
}

func DoDigestJob(ctx context.Context, dbClient *db.Client, mailClient MailClient) (err error) {
	log.InfoContext(ctx, "Starting mail send job")
	userIds, err := GetUsers(ctx, dbClient)
	if err != nil {
		return
	}

	wg := new(sync.WaitGroup)
	for _, user := range userIds {
		wg.Add(1)
		go func(user User) {
			defer wg.Done()
			message, err := GetEmailForUser(ctx, dbClient, user)
			err = mailClient.SendMail(ctx, &user, message)
			if err != nil {
				log.ErrContext(ctx, "Failed sending mail", err)
			}
		}(user)
	}
	wg.Wait()
	log.InfoContext(ctx, "Mail send job done")
	return nil
}
