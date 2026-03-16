package packets

import (
	"bufio"
	"errors"

	"github.com/maekilae/gophermc/internal/protocol/types"
)

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      uint16
	NextState       int32 // 1 for Status, 2 for Login
}

func (h *Handshake) ID() int32 {
	return 0x00
}

func (h *Handshake) Read(r *bufio.Reader) error {
	var version types.VarInt
	var address types.StringN
	var port types.UnsignedShort
	var state types.VarInt

	if err := version.Read(r); err != nil {
		return err
	}
	if err := address.Read(r); err != nil {
		return err
	}
	if err := port.Read(r); err != nil {
		return err
	}
	if err := state.Read(r); err != nil {
		return err
	}
	h.ProtocolVersion = int32(version)
	h.ServerAddress = string(address)
	h.ServerPort = uint16(port)
	h.NextState = int32(state)

	return nil
}

func (h *Handshake) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
