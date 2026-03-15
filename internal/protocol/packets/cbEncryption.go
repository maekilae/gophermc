package packets

import (
	"bufio"
	"errors"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type EncryptionRequest struct {
	ServerID   types.StringN
	PubKey     types.ByteArray
	Token      types.ByteArray
	ShouldAuth types.Boolean
}

func (pk *EncryptionRequest) ID() int32 {
	return 0x01
}

func (pk *EncryptionRequest) Read(w *bufio.Reader) error {
	return errors.New("Not implemented")
}

func (en *EncryptionRequest) Write(w *bufio.Writer) error {
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
