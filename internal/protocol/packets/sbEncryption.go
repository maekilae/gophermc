package packets

import (
	"bufio"
	"errors"

	"github.com/maekilae/gophermc/internal/protocol/types"
)

type EncryptionResponse struct {
	SharedSecret types.ByteArray
	VerifyToken  types.ByteArray
}

func (pk *EncryptionResponse) ID() int32 {
	return 0x01
}

func (pk *EncryptionResponse) Read(w *bufio.Reader) error {
	err := pk.SharedSecret.Read(w)
	if err != nil {
		return err
	}
	err = pk.VerifyToken.Read(w)
	if err != nil {
		return err
	}
	return nil
}

func (pk *EncryptionResponse) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
