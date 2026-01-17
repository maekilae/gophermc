package client

import (
	"net"

	"codeberg.org/makila/minecraftgo/internal/packet"
)

type Client struct {
	username string
	uuid     packet.UUID

	addr    string
	tcpaddr *net.TCPAddr
}
