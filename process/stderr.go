package process

import (
	"io"
	"os"
)

// StderrReadCloser handles process's Stdout
//   this reader also exports its outputs to stdout
type StderrReadCloser struct {
	readCloser io.ReadCloser
}

// NewStderr make stdin read and closer
func NewStderr(r io.ReadCloser) *StderrReadCloser {
	return &StderrReadCloser{r}
}

func (s StderrReadCloser) Read(p []byte) (n int, err error) {
	n, err = s.readCloser.Read(p)
	os.Stderr.Write(p[:n])

	return
}

// Close ...
func (s StderrReadCloser) Close() error {
	return s.readCloser.Close()
}
