package urls

import "strings"

const (
	Root   = "/"
	Assets = "/assets"
	Login  = "/login"
	App    = "/app"
	Page   = "/app/page"
)

func GroupURL(prefix, url string) string {
	return strings.TrimPrefix(url, prefix)
}
