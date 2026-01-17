package packet

import (
	"bytes"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Encryption struct {
	ServerID   string
	PubKey     []byte
	Token      []byte
	ShouldAuth bool
}

func (en *Encryption) ID() int32 {
	return 0x01
}

func (en *Encryption) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(types.WriteVarInt(len(en.ServerID)))
	buf.Write(en.PubKey)
	buf.Write(en.Token)
	if en.ShouldAuth {
		buf.WriteByte(0x01)
	} else {
		buf.WriteByte(0x00)
	}
	return buf.Bytes(), nil

}

