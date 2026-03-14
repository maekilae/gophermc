package player

import (
	"github.com/google/uuid"
)

type Properties struct {
	Name      string
	Value     string
	Signature string
}

type Player struct {
	Username string
	Uuid     uuid.UUID
	Props    Properties
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
