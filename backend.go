package backend

import (
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

type Opener interface {
	Open(string) (io.ReadCloser, error)
}

type Backend interface {
	Opener
	/*
		Ls(path string) ([]os.FileInfo, error)
		Mkdir(path string) error
		Rm(path string) error
		Mv(from string, to string) error
		Save(path string, file io.Reader) error
		Touch(path string) error
	*/
}

var reURL = regexp.MustCompile(`(?i)^[a-z][a-z0-9]*://`)

func Parse(uri string) (*url.URL, Backend, error) {
	if !reURL.MatchString(uri) {
		uri = "file:///" + uri
	}

	url, err := url.Parse(uri)
	if err != nil {
		return nil, nil, err
	}
	switch scheme := strings.ToLower(url.Scheme); scheme {
	case "file":
		return url, &FileBackend{}, nil
	case "http", "https":
		return url, &HTTPBackend{}, nil
	default:
		return nil, nil, errors.New("unsupported URL scheme: " + scheme + ": " + uri)
	}
}

type AutoBackend struct{}

func Default() Backend {
	return &AutoBackend{}
}

func (_ *AutoBackend) Open(uri string) (io.ReadCloser, error) {
	url, b, err := Parse(uri)
	if err != nil {
		return nil, err
	}

	switch url.Scheme {
	case "file":
		return b.Open(url.Path)
	case "http", "https":
		return b.Open(url.String())
	}
	return nil, ErrNotImplemented
}

func (b *AutoBackend) ReadFile(uri string) ([]byte, error) {
	r, err := b.Open(uri)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

var ErrNotImplemented = errors.New("not implemented")
