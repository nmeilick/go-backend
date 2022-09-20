package backend

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPBackend allows accessing files using the http/https protocol
type HTTPBackend struct {
	client *http.Client
}

func (b *HTTPBackend) Ensure() *HTTPBackend {
	if b != nil {
		return b
	}
	return &HTTPBackend{}
}

func (b *HTTPBackend) Client() *http.Client {
	b = b.Ensure()
	if b.client == nil {
		b.client = http.DefaultClient
	}
	return b.client
}

type HTTPBody struct {
	r io.ReadCloser
}

func (b *HTTPBody) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

func (b *HTTPBody) Close() error {
	io.Copy(ioutil.Discard, b.r)
	return b.r.Close()
}

func (b *HTTPBackend) Open(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.Client().Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, NewStatusError(resp.StatusCode)
	}
	return &HTTPBody{r: resp.Body}, nil
}

type StatusError struct {
	code int
}

func (e *StatusError) Error() string {
	if e == nil {
		return fmt.Sprintf("%d (%s)", http.StatusOK, http.StatusText(http.StatusOK))
	}
	return fmt.Sprintf("%d (%s)", e.code, http.StatusText(e.code))
}

func NewStatusError(code int) *StatusError {
	return &StatusError{code: code}
}
