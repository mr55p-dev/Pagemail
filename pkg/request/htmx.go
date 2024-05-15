package request

import "net/http"

func IsHtmx(r *http.Request) bool {
	return r.Header.Get("Hx-Request") == "true"
}
