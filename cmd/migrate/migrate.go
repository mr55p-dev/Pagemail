package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	migrations "github.com/mr55p-dev/pagemail/db"
)

func main() {
	u, err := url.Parse(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	db := dbmate.New(u)
	db.FS = migrations.FS

	migrations, err := db.FindMigrations()
	if err != nil {
		panic(err)
	}

	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}

	fmt.Println("\nApplying...")
	err = db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
}
