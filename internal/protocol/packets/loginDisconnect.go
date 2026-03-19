package packets

import (
	"bufio"
	"encoding/json"
	"errors"

	"github.com/maekilae/gophermc/internal/protocol/types"
)

type LoginDisconnect struct {
	Reason types.StringN `json:"text"`
}

func (pk LoginDisconnect) ID() int32 {
	return 0x00
}

func (pk LoginDisconnect) Read(r *bufio.Reader) error {
	return errors.New("Not implemented")
}

func (pk LoginDisconnect) Write(w *bufio.Writer) error {
	jsonBytes, err := json.Marshal(pk)
	if err != nil {
		return err
	}
	length := types.VarInt(len(jsonBytes))
	err = length.Write(w)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}
