package adding

import (
	"io"
)

type Writer interface {
	Write(r io.Reader, filepath string) error
}

type Adder struct {
	writer Writer
}

func NewAdder(writer Writer) *Adder {
	return &Adder{writer: writer}
}

func (a *Adder) Add(r io.Reader, filepath string) (string, error) {
	err := a.writer.Write(r, filepath)
	return filepath, err
}
