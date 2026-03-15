package packets

import (
	"bufio"
	"errors"
)

type LoginAcknowledge struct {
}

func (s *LoginAcknowledge) ID() int32 {
	return 0x03
}

func (s *LoginAcknowledge) Read(w *bufio.Reader) error {
	return nil
}

func (s *LoginAcknowledge) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
