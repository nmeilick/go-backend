package backend

import (
	"io"
	"os"
	"path/filepath"
)

// FileBackend allows accessing locally mounted filesystems
type FileBackend struct{}

func (b *FileBackend) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (b *FileBackend) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
