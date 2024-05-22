package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mr55p-dev/pagemail/db"
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/tools"
)

var recipient = flag.String("to", "", "Address to send the test email to")

func loadUserWithPages(ctx context.Context, db *sql.DB, address string, now time.Time) error {
	uid := tools.GenerateNewId(5)
	_, err := queries.New(db).CreateUser(ctx, queries.CreateUserParams{
		ID:         uid,
		Username:   "test",
		Email:      address,
		Subscribed: true,
	})
	if err != nil {
		return err
	}

	_, err = queries.New(db).CreatePage(ctx, queries.CreatePageParams{
		ID:     tools.GenerateNewId(5),
		UserID: uid,
		Url:    "https://www.google.com",
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
	client := mail.NewAwsSender(ctx, cfg)

	//
	now := time.Now()
	err = loadUserWithPages(ctx, conn, *recipient, now)
	if err != nil {
		panic(err)
	}

	err = mail.DigestJob(ctx, conn, client, now)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job complete")
}
