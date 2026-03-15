package packets

import (
	"bufio"
	"errors"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type LoginSuccess struct {
	GameProfile types.GameProfile
}

func (pk *LoginSuccess) ID() int32 {
	return 0x02
}

func (pk *LoginSuccess) Read(w *bufio.Reader) error {
	return errors.New("Not implemented")
}

func (pk *LoginSuccess) Write(w *bufio.Writer) error {
	err := pk.GameProfile.Write(w)
	if err != nil {
		return nil
	}
	return nil
}
