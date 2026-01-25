package packet

import (
	"encoding/json"
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

func (pk *ServerStatus) ID() int32 {
	return 0x00
}

func (pk *ServerStatus) Marshal() ([]byte, error) {
	return json.Marshal(pk)
}
