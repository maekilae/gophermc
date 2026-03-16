package player

type PlayerData struct {
	UUID       string     `json:"id"`
	Username   string     `json:"name"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
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
