package request

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mr55p-dev/pagemail/internal/logging"
	"github.com/stretchr/testify/assert"
)

type BindTestStruct struct {
	FormString  string `form:"form-string"`
	FormInt     int    `form:"form-int"`
	FormBool    bool   `form:"form-bool"`
	QueryString string `query:"query-string"`
	QueryInt    int    `query:"query-int"`
	QueryBool   bool   `query:"query-bool"`
}

func TestBind(t *testing.T) {
	logging.SetHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	assert := assert.New(t)
	req := httptest.NewRequest(
		http.MethodPost,
		"/?query-string=world&query-int=456&query-bool=true",
		strings.NewReader("form-string=hello&form-int=123&form-bool=false"),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	bound := new(BindTestStruct)
	err := Bind(bound, req)
	assert.NoError(err)
	assert.Equal(&BindTestStruct{
		FormString:  "hello",
		FormInt:     123,
		FormBool:    false,
		QueryString: "world",
		QueryInt:    456,
		QueryBool:   true,
	}, bound)
}
