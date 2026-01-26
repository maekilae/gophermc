package clientbound

import (
	"encoding/json"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type ServerStatus struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favicon"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Player `json:"sample"`
}

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Description struct {
	Text string `json:"text"`
}

func (pk ServerStatus) ID() int32 {
	return 0x00
}

func (pk ServerStatus) Marshal() (buf []byte, err error) {
	data, _ := json.Marshal(pk)
	types.VarInt(len(data)).ToBytes(buf)
	buf = append(buf,data...)
	return
}
