package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/render"
	"github.com/mr55p-dev/pagemail/internal/timer"
)

type MailSender interface {
	Send(context.Context, string, io.Reader) error
}

type MailDbReader interface {
	ReadUsersWithMail(context.Context) ([]db.User, error)
	ReadPagesByUserBetween(context.Context, string, time.Time, time.Time) ([]db.Page, error)
}

func MailGo(ctx context.Context, reader MailDbReader, sender MailSender) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.Local)
	timer := timer.NewCronTimer(time.Hour*24, start)
	defer timer.Stop()

	for now := range timer.T {
		ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
		err := MailJob(ctx, reader, sender, now)
		cancel()
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to send digest", "error", err.Error())
		}
	}
}

func MailJob(ctx context.Context, reader MailDbReader, sender MailSender, now time.Time) error {
	// Get users with mail enabled
	users, err := reader.ReadUsersWithMail(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get user list: %w", err)
	}

	// Dispatch jobs to the other goroutines
	var errCount int32
	jobs := make(chan db.User)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range jobs {
				err := sendMailToUser(ctx, &user, sender, reader, now)
				if err != nil {
					errChan <- err
				}
			}
		}()
	}

	// Read out all the errors
	errList := make([]error, 0)
	go func() {
		for err := range errChan {
			atomic.AddInt32(&errCount, 1)
			errList = append(errList, err)
		}
	}()

	// Dispatch all the jobs
	for _, user := range users {
		jobs <- user
	}
	close(jobs)
	wg.Wait()

	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}

func sendMailToUser(ctx context.Context, user *db.User, sender MailSender, db MailDbReader, now time.Time) error {
	logger := logger.With("user", user.Email)
	logger.DebugCtx(ctx, "Generating mail for user")

	// Read the users pages
	start := now.Add(-time.Hour * 24)
	end := now
	pages, err := db.ReadPagesByUserBetween(ctx, user.Id, start, end)
	if err != nil {
		return err
	}

	// Brea
	if len(pages) == 0 {
		logger.DebugCtx(ctx, "No pages found")
		return nil
	}

	// Generate the mail and send it
	buf := bytes.Buffer{}
	err = render.MailDigest(now, user.Name, pages).Render(ctx, &buf)
	if err != nil {
		return err
	}
	err = sender.Send(ctx, user.Email, &buf)
	if err != nil {
		return err
	}
	return nil
}
