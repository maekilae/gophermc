package packets

import (
	"bufio"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Encryption struct {
	ServerID   types.StringN
	PubKey     types.ByteArray
	Token      types.ByteArray
	ShouldAuth types.Boolean
}

func (en Encryption) ID() int32 {
	return 0x01
}

func (en Encryption) Write(w *bufio.Writer) error {
	// buf := new(bytes.Buffer)

	err := en.ServerID.Write(w)
	if err != nil {
		return err
	}
	err = en.PubKey.Write(w)
	if err != nil {
		return err
	}
	err = en.Token.Write(w)
	if err != nil {
		return err
	}
	err = en.ShouldAuth.Write(w)
	if err != nil {
		return err
	}
	return nil

}
