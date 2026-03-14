package packets

import (
	"bufio"
	"errors"
)

type StatusRequest struct {
}

func (s *StatusRequest) ID() int32 {
	return 0x00
}

func (s *StatusRequest) Read(w *bufio.Reader) error {
	return nil
}

func (s *StatusRequest) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
