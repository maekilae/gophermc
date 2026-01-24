package player

import (
	"bytes"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
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

func (p *Player) Marshal() []byte {
	buf := new(bytes.Buffer)

	uuid := types.WriteUUID(p.Uuid)
	buf.Write(uuid)
	buf.Write(types.WriteString(p.Username))
	buf.Write(types.WriteString(p.Props.Name))
	buf.Write(types.WriteString(p.Props.Value))
	buf.Write(types.WriteString(p.Props.Signature))
	return buf.Bytes()

}
