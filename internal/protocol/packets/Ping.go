package packets

import (
	"bufio"

	"github.com/maekilae/gophermc/internal/protocol/types"
)

type Ping struct {
	timestamp types.SignedLong
}

func (pk *Ping) ID() int32 {
	return 0x01
}

func (pk *Ping) Read(w *bufio.Reader) error {
	return pk.timestamp.Read(w)
}

func (pk *Ping) Write(w *bufio.Writer) error {
	return pk.timestamp.Write(w)
}
