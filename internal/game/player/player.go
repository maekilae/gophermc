package player

import (
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Properties struct {
	Name      types.StringN
	Value     types.StringN
	Signature types.StringN
}

type Player struct {
	Username   types.StringN
	Uuid       types.UUID
	Properties Properties
}

// func (p *Player) Marshal() []byte {
// 	buf := new(bytes.Buffer)

// 	uuid := types.WriteUUID(p.Uuid)
// 	buf.Write(uuid)
// 	types.StringN(p.Username).ToBytes(buf)
// 	buf.Write(types.StringN(p.Username))
// 	buf.Write(types.WriteString(p.Props.Name))
// 	buf.Write(types.WriteString(p.Props.Value))
// 	buf.Write(types.WriteString(p.Props.Signature))
// 	return buf.Bytes()

// }
