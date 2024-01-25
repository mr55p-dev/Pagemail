package main

import (
	"flag"
	"fmt"

	"github.com/mr55p-dev/pagemail/pkg/tools"
)

func main() {
	tkn := tools.GenerateNewShortcutToken()
	silent := flag.Bool("silent", false, "Suppress all output except for the token")
	flag.Parse()
	if silent == nil || !(*silent) {
		fmt.Printf("Generated new token: %v\n", tkn)
	} else {
		fmt.Println(tkn)
	}
}
