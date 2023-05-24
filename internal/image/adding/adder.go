package adding

import (
	"io"

	"github.com/Eviljeks/blurer/internal/storage"
)

type Adder struct {
	storage *storage.Storage
}

func NewAdder(storage *storage.Storage) *Adder {
	return &Adder{storage: storage}
}

func (a *Adder) Add(r io.Reader, filepath string) (string, error) {
	err := a.storage.Write(r, filepath)
	return filepath, err
}
