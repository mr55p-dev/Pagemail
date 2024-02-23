package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/mr55p-dev/pagemail/internal/auth"
	"github.com/mr55p-dev/pagemail/internal/db"
	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/mr55p-dev/pagemail/internal/mail"
	"github.com/mr55p-dev/pagemail/internal/tools"
	"github.com/urfave/cli/v2"
)

var (
	dbClient   *db.Client
	authClient auth.Authorizer
)

func init() {
	authClient = auth.NewSecureAuthorizer(context.TODO())
}

func GetUser(ctx *cli.Context) (user *db.User, err error) {
	if id := ctx.String("user-id"); id != "" {
		user, err = dbClient.ReadUserById(context.TODO(), id)
	} else if email := ctx.String("user-email"); email != "" {
		user, err = dbClient.ReadUserByEmail(context.TODO(), email)
	} else {
		return nil, fmt.Errorf("No user parameters given")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func main() {
	log := logging.NewVoid()
	app := cli.App{
		Name:        "pmtk",
		Description: "Toolkit for pagemail services and debugging",
		Suggest:     true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "db",
				EnvVars:  []string{"PM_DB_PATH"},
				Required: true,
				Action: func(ctx *cli.Context, p cli.Path) error {
					dbClient = db.NewClient(p, log)
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "user-id",
					},
					&cli.StringFlag{
						Name: "user-email",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name: "get",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "json",
							},
						},
						Action: func(ctx *cli.Context) (err error) {
							user, err := GetUser(ctx)
							if err != nil {
								return err
							}
							if ctx.Bool("json") {
								var data []byte
								data, err = json.Marshal(user)
								if err != nil {
									return
								}
								os.Stdout.Write(data)
							} else {
								fmt.Print(user)
							}
							return
						},
					},
					{
						Name: "set-password",
						Action: func(ctx *cli.Context) error {
							user, err := GetUser(ctx)
							if err != nil {
								return err
							}
							password := ctx.Args().First()
							if password == "" {
								return fmt.Errorf("Password must be provided")
							}
							user.Password = authClient.GenPasswordHash(password)
							return dbClient.UpdateUser(
								context.TODO(),
								user,
							)
						},
					},
				},
			},
			{
				Name: "mail",
				Subcommands: []*cli.Command{
					{
						Name: "run-all",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "real-client",
								Aliases: []string{"real"},
								Usage:   "Use the real ses client",
							},
						},
						Action: func(ctx *cli.Context) error {
							var mailClient mail.MailClient
							if ctx.Bool("real-client") {
								mailClient = mail.NewSesMailClient(context.TODO())
							} else {
								mailClient = &mail.TestClient{}
							}
							slog.Info("Starting to send mail")
							mail.DoDigestJob(
								context.TODO(),
								dbClient,
								mailClient,
							)
							return nil
						},
					},
					{
						Name: "generate",
						Args: false,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "user-id",
								Aliases:  []string{"uid", "id"},
								EnvVars:  []string{"PM_TEST_USER"},
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							id := ctx.String("user-id")
							user, err := dbClient.ReadUserById(context.TODO(), id)
							if err != nil {
								return err
							}
							msg, err := mail.GetEmailForUser(
								context.TODO(),
								dbClient,
								mail.User{
									Id:    id,
									Name:  user.Name,
									Email: user.Email,
								},
							)
							if err != nil {
								return err
							}
							fmt.Print(msg)
							return nil
						},
					},
				},
			},
			{
				Name: "generate",
				Subcommands: []*cli.Command{
					{
						Name: "shortcut-token",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "silent",
								Aliases: []string{"s"},
							},
						},
						Action: func(ctx *cli.Context) error {
							tkn := tools.GenerateNewShortcutToken()
							silent := ctx.Bool("silent")
							if silent {
								fmt.Println(tkn)
							} else {
								slog.Info("Generated new shortcut token", "value", tkn)
							}
							return nil
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
