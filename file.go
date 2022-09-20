package backend

import (
	"io"
	"os"
)

// FileBackend allows accessing locally mounted filesystems
type FileBackend struct{}

func (b *FileBackend) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
