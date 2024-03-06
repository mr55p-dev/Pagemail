package httpit

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	Field1 string `form:"field1"`
	Field2 string `query:"field2"`
}

func TestBind(t *testing.T) {
	assert := assert.New(t)
	r, err := http.NewRequest("POST", "/?field2=field2", strings.NewReader("field1=field1"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", "13")
	r.ParseForm()
	fmt.Printf("r.Form.Encode(): %v\n", r.Form)
	assert.Nil(err)
	assert.Equal("field1", r.Form.Get("field1"))
	assert.Equal("field2", r.URL.Query().Get("field2"))

	out := new(Example)
	err = bind(out, r)
	assert.Nil(err)
	assert.Equal("field1", out.Field1)
	assert.Equal("field2", out.Field2)
}
