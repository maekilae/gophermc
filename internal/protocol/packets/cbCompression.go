package packets

import (
	"bufio"
	"errors"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Compression struct {
	Threshold types.VarInt
}

func (pk *Compression) ID() int32 {
	return 0x03
}

func (pk *Compression) Read(w *bufio.Reader) error {
	return errors.New("Not implemented")
}

func (pk *Compression) Write(w *bufio.Writer) error {
	err := pk.Threshold.Write(w)
	if err != nil {
		return err
	}
	return nil
}
