package configLoader

import (
	"fmt"
	"strings"
)

func getEnvName(key string, prefix string) string {
	replacer := strings.NewReplacer(
		"-", "_",
		".", "_",
	)
	if prefix != "" {
		prefix = strings.ToUpper(prefix)
		prefix = fmt.Sprintf("%s_", prefix)
	}
	return fmt.Sprintf("%s%s", prefix, strings.ToUpper(replacer.Replace(key)))
}
