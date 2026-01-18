package player

import (
	"net"

	"codeberg.org/makila/minecraftgo/internal/protocol/packet"
)

type Player struct {
	Name string
	Uuid string

	Conn  net.Conn
	State packet.State
}
