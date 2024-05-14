package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePassword(t *testing.T) {
	assert := assert.New(t)
	rdr := strings.NewReader("password\r")
	pass, err := parsePassword(rdr)
	assert.NoError(err)
	assert.Equal([]byte("password"), pass)
}
