package packets

import (
	"bufio"
	"errors"

	"github.com/maekilae/gophermc/internal/protocol/types"
)

type LoginStart struct {
	Username types.StringN
	UUID     types.UUID
}

func (pk *LoginStart) ID() int32 {
	return 0x00
}

func (pk *LoginStart) Read(w *bufio.Reader) error {
	err := pk.Username.Read(w)
	if err != nil {
		return err
	}
	err = pk.UUID.Read(w)
	if err != nil {
		return err
	}
	return nil
}

func (pk *LoginStart) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
