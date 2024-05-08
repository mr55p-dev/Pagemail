package main

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
	QueryString string `query:"query-string"`
	QueryInt    int    `query:"query-int"`
}

func TestBind(t *testing.T) {
	logging.SetHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	assert := assert.New(t)
	req := httptest.NewRequest(
		http.MethodPost,
		"/?query-string=world&query-int=456",
		strings.NewReader("form-string=hello&form-int=123"),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	bound := new(BindTestStruct)
	err := Bind(bound, req)
	assert.NoError(err)
	assert.Equal(&BindTestStruct{
		FormString:  "hello",
		FormInt:     123,
		QueryString: "world",
		QueryInt:    456,
	}, bound)
}
