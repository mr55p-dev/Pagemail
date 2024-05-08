package mail

import (
	"context"
	"github.com/mr55p-dev/pagemail/internal/db"
	"sync"
	"time"
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
	logger.DebugCtx(ctx, "Looking up pages")
	since := Yesterday()

	pages, err := dbClient.ReadPagesByUserId(ctx, user.Id, -1)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed looking up pages", err)
	}
	pages = FilterPages(ctx, pages, since)
	logger.DebugCtx(ctx, "Found pages", "count", len(pages))
	messageBytes, err := GenerateMailBody(ctx, &user, pages, since)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed while generating mail body", err)
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

func DoDigestJob(ctx context.Context, dbClient *db.Client, mailClient Sender) (err error) {
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
			err = mailClient.Send(ctx, &user, message)
			if err != nil {
				logger.ErrorCtx(ctx, "Failed sending mail", err)
			}
		}(user)
	}
	wg.Wait()
	logger.InfoCtx(ctx, "Mail send job done")
	return nil
}
