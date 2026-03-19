package packets

import (
	"bufio"
	"errors"
)

type Finish struct {
}

func (*Finish) ID() uint32 {
	return 0x03
}

func (pk *Finish) Read(w *bufio.Reader) error {
	return errors.New("Not implemented")
}

func (pk *Finish) Write(w *bufio.Writer) error {
	return nil
}

type FinishAcknowledge struct {
}

func (*FinishAcknowledge) ID() uint32 {
	return 0x03
}

func (pk *FinishAcknowledge) Read(w *bufio.Reader) error {
	return nil
}

func (pk *FinishAcknowledge) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
