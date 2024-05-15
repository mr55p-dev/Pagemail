package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodedReader(t *testing.T) {
	assert := assert.New(t)
	baseStream := strings.NewReader("hello, world!")
	encodedDest := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, encodedDest)
	wrapper := &EncodedReader{
		Source:     baseStream,
		encoderIn:  encoder,
		encoderOut: encodedDest,
	}

	res, err := io.ReadAll(wrapper)
	assert.NoError(err)
	assert.Equal("aGVsbG8sIHdvcmxkIQ==", string(res))
}

func TestBase64EncodedReader(t *testing.T) {
	src := strings.NewReader("this is a test")
	reader := NewBase64EncodedReader(src, base64.StdEncoding)
	res, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, []byte("dGhpcyBpcyBhIHRlc3Q="), res)
}
