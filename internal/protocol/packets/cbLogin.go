package packets

import (
	"bytes"

	"codeberg.org/makila/minecraftgo/internal/game/player"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type LoginSuccess struct {
	Profile player.Player
}

func (ls LoginSuccess) ID() int32 {
	return 0x02
}
func (ls LoginSuccess) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	// buf.Write(ls.Profile.Marshal())

	return buf.Bytes(), nil
}

type SetCompression struct {
	Threshold types.VarInt
}

func (pk SetCompression) ID() int32 {
	return 0x03
}

func (pk SetCompression) Marshal() (buf []byte, err error) {
	// pk.Threshold.ToBytes(buf)
	return
}
