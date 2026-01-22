package packet

import (
	"bytes"

	"codeberg.org/makila/minecraftgo/internal/game/player"
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
	buf.Write(types.WriteString(en.ServerID))
	buf.Write(types.WriteByteArray(en.PubKey))
	buf.Write(types.WriteByteArray(en.Token))
	if en.ShouldAuth {
		buf.WriteByte(0x01)
	} else {
		buf.WriteByte(0x00)
	}
	return buf.Bytes(), nil

}

type LoginSuccess struct {
	Profile player.Player
}

func (ls *LoginSuccess) ID() int32 {
	return 0x02
}
func (ls *LoginSuccess) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	uuid := types.WriteUUID(ls.Profile.Uuid)
	buf.Write(uuid)
	buf.Write(types.WriteString(ls.Profile.Username))
	buf.Write(types.WriteString(ls.Profile.Props.Name))
	buf.Write(types.WriteString(ls.Profile.Props.Value))
	buf.Write(types.WriteString(ls.Profile.Props.Signature))

	return buf.Bytes(), nil
}
