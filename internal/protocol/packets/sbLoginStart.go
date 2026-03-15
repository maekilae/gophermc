package packets

import (
	"bufio"
	"errors"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type LoginStart struct {
	Username types.StringN
	UUID     types.UUID
}

func (pk *LoginStart) ID() int32 {
	return 0x01
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
