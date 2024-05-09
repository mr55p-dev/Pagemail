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

type MailNoOp struct{}

func (MailNoOp) Send(context.Context, string, io.Reader) error {
	logger.Debug("No-op mail send")
	return nil
}

func MailGo(ctx context.Context, router *Router) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.Local)
	timer := timer.NewCronTimer(time.Hour*24, start)
	defer timer.Stop()

	for now := range timer.T {
		ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
		err := router.MailJob(ctx, now)
		cancel()
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to send digest", "error", err.Error())
		}
	}
}

func (router *Router) MailJob(ctx context.Context, now time.Time) error {
	// Get users with mail enabled
	users, err := router.DBClient.ReadUsersWithMail(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get user list: %w", err)
	}

	// Dispatch jobs to the other goroutines
	var errCount int32
	jobs := make(chan db.User)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doMailJob(
				ctx,
				jobs,
				errChan,
				router.MailClient,
				router.DBClient,
				now,
			)
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

func doMailJob(ctx context.Context, job chan db.User, errors chan error, sender MailSender, db *db.Client, now time.Time) {
	for user := range job {
		logger := logger.With("user", user.Email)
		logger.DebugCtx(ctx, "Generating mail for user")

		start := now.Add(-time.Hour * 24)
		end := now
		pages, err := db.ReadPagesByUserBetween(ctx, user.Id, start, end)
		if err != nil {
			errors <- err
			continue
		}

		if len(pages) == 0 {
			logger.DebugCtx(ctx, "No pages for user")
			continue
		}

		// Generate the mail and send it
		buf := bytes.Buffer{}
		err = render.MailDigest(now, user.Name, pages).Render(ctx, &buf)
		if err != nil {
			errors <- err
		}
		err = sender.Send(ctx, user.Email, &buf)
		if err != nil {
			errors <- err
		}
	}
}
