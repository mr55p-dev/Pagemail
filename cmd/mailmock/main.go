package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

var recipient = flag.String("to", "", "Address to send the test email to")

func loadUserWithPages(ctx context.Context, queries *dbqueries.Queries, address string, now time.Time) error {
	uid := tools.GenerateNewId(5)
	err := queries.CreateUser(ctx, dbqueries.CreateUserParams{
		ID:         uid,
		Username:   "test",
		Password:   []byte("password"),
		Email:      address,
		Subscribed: true,
		Created:    now,
		Updated:    now,
	})
	if err != nil {
		return err
	}

	err = queries.CreatePage(ctx, dbqueries.CreatePageParams{
		ID:      tools.GenerateNewId(5),
		UserID:  uid,
		Url:     "https://www.google.com",
		Created: now.Add(-time.Hour),
		Updated: now.Add(-time.Hour),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "mailmock [options] [message body]\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	// Setup the database
	conn := db.MustConnect(ctx, ":memory:")
	db.MustLoadSchema(ctx, conn)
	queries := dbqueries.New(conn)
	client := mail.NewSesMailClient(ctx, cfg)

	//
	now := time.Now()
	err = loadUserWithPages(ctx, queries, *recipient, now)
	if err != nil {
		panic(err)
	}

	err = mail.MailJob(ctx, queries, client, now)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job complete")
}
