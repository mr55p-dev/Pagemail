package response

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/mr55p-dev/pagemail/internal/pmerror"
	"github.com/mr55p-dev/pagemail/internal/render/components/popup"
	"github.com/mr55p-dev/pagemail/pkg/request"
)

type MessageFn = func(string, string) templ.Component

var ErrorComponent MessageFn = func(s1, s2 string) templ.Component {
	return popup.New(s1, popup.WithBody(s2))
}
var OkComponent MessageFn = func(s1, s2 string) templ.Component {
	return popup.New(s1, popup.WithBody(s2))
}

func Component(component templ.Component, w http.ResponseWriter, r *http.Request) {
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering response", http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	data, ok := err.(*pmerror.PMError)
	if !ok {
		data = pmerror.ErrUnspecified
	}
	if request.IsHtmx(r) {
		w.WriteHeader(data.Status)
		Component(ErrorComponent("Error", data.Message), w, r)
	} else {
		http.Error(w, fmt.Sprintf("Error: %s", data.Message), data.Status)
	}
}

func Success(message string, w http.ResponseWriter, r *http.Request) {
	if request.IsHtmx(r) {
		Component(OkComponent("Success", message), w, r)
	} else {
		_, err := fmt.Fprint(w, message)
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
		}
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	if request.IsHtmx(r) {
		w.Header().Add("HX-Redirect", url)
	} else {
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusSeeOther)
	}
}
