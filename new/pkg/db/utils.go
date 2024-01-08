package db

import "strings"

func GetPlaceholder(values []string) string {
	ps := make([]string, len(values))
	for i := 0; i < len(ps); i++ {
		ps[i] = "?"
	}
	return strings.Join(ps, ",")
}
