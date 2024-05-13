package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mr55p-dev/gonk"
	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/dbqueries"
)

var userId = flag.String("user-id", "", "User ID")

type Config struct {
	Url string `config:"db_url"`
}

func main() {
	ctx := context.Background()
	usage := func() {
		fmt.Fprintf(os.Stderr, "passreset can reset passwords in the database")
		flag.PrintDefaults()
	}
	flag.Usage = usage
	flag.Parse()

	cfg := new(Config)
	err := gonk.LoadConfig(cfg, gonk.EnvironmentLoader(""))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v", err)
		os.Exit(1)
	}

	if *userId == "" {
		fmt.Fprintf(os.Stderr, "user_id is required")
		usage()
		os.Exit(1)
	}

	password := new(bytes.Buffer)
	io.Copy(os.Stdin, password)
	passwordHash := auth.HashPassword(password.Bytes())

	conn := dbqueries.MustGetDB(ctx, cfg.Url)
	queries := dbqueries.New(conn)
	err = queries.UpdateUserPassword(ctx, dbqueries.UpdateUserPasswordParams{
		Password: passwordHash,
		ID:       *userId,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating password: %v", err)
		os.Exit(1)
	}
	fmt.Println("Password updated sucessfully")
	return
}
