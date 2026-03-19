package packets

import (
	"bufio"
	"encoding/json"
	"errors"

	"github.com/maekilae/gophermc/internal/protocol/types"
)

type StatusResponse struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favicon,omitempty"` // omitempty is good here if you don't have an icon
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Player `json:"sample,omitempty"`
}

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Description struct {
	Text string `json:"text"`
}

func (pk StatusResponse) ID() int32 {
	return 0x00
}

func (pk StatusResponse) Read(r *bufio.Reader) error {
	return errors.New("Not implemented")
}

func (pk StatusResponse) Write(w *bufio.Writer) error {
	jsonBytes, err := json.Marshal(pk)
	if err != nil {
		return err
	}
	length := types.VarInt(len(jsonBytes))
	err = length.Write(w)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

type StatusRequest struct {
}

func (s *StatusRequest) ID() int32 {
	return 0x00
}

func (s *StatusRequest) Read(w *bufio.Reader) error {
	return nil
}

func (s *StatusRequest) Write(w *bufio.Writer) error {
	return errors.New("Not implemented")
}
