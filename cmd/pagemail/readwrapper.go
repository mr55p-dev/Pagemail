package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
)

type EncodedReader struct {
	Source     io.Reader
	encoderIn  io.WriteCloser
	encoderOut io.Reader
}

func NewBase64EncodedReader(source io.Reader, encoding *base64.Encoding) *EncodedReader {
	dest := new(bytes.Buffer)
	encoder := base64.NewEncoder(encoding, dest)

	return &EncodedReader{
		Source:     source,
		encoderIn:  encoder,
		encoderOut: dest,
	}
}

func (r *EncodedReader) Read(to []byte) (int, error) {
	eof := false
	mid := make([]byte, len(to))
	n, err := r.Source.Read(mid)
	if err != nil {
		if errors.Is(err, io.EOF) {
			eof = true
		} else {
			return 0, err
		}
	}

	if eof {
		_ = r.encoderIn.Close()
	} else {
		_, err = r.encoderIn.Write(mid[:n])
		if err != nil {
			return 0, err
		}
	}

	n, err = r.encoderOut.Read(to)
	if err != nil {
		return n, err
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}
