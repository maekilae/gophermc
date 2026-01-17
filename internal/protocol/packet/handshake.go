package packet

import (
	"bufio"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
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

func (h *Handshake) Unmarshal(r *bufio.Reader) error {
	var e error
	h.ProtocolVersion, e = types.ReadVarInt(r)
	if e != nil {
		return e
	}

	h.ServerAddress, e = types.ReadString(r)
	if e != nil {
		return e
	}

	h.ServerPort, e = types.ReadUShort(r)
	if e != nil {
		return e
	}

	h.NextState, e = types.ReadVarInt(r)
	if e != nil {
		return e
	}

	return nil
}
