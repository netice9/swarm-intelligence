package api

import (
	"errors"
	"io"
	"net/http"
)

type flushWriter struct {
	io.Writer
}

func (fw flushWriter) Write(data []byte) (int, error) {
	flusher, ok := fw.Writer.(http.Flusher)
	if !ok {
		return 0, errors.New("Streaming is not supported")
	}
	n, err := fw.Writer.Write(data)
	if err != nil {
		return n, err
	}
	flusher.Flush()
	return n, err
}
