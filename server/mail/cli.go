package mail

import (
	"log"
	"pagemail/server/models"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func MailCommand(app *pocketbase.PocketBase, cfg *models.PMContext) *cobra.Command {
	cmd := func(c *cobra.Command, args []string) {
		err := Mailer(app, cfg)
		if err != nil {
			log.Printf("Mailer failed with error %s", err)
		} else {
			log.Print("Mailer succeeded")
		}
		return
	}

	return &cobra.Command{
		Use: "mail-all",
		Run: cmd,
	}
}
