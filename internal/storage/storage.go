package storage

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	pathPrefix string
}

func NewStorage(pathPrefix string) *Storage {
	if pathPrefix[len(pathPrefix)-1] != filepath.Separator {
		pathPrefix += string(filepath.Separator)
	}
	return &Storage{pathPrefix: pathPrefix}
}

func (s *Storage) Open(src string) (io.ReadWriteCloser, error) {
	output, err := os.OpenFile(s.pathPrefix+src, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *Storage) Write(r io.Reader, filepath string) error {
	f, err := s.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)

	return err
}
