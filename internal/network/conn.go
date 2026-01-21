package network

import (
	"bufio"
	"io"
	"net"

	"codeberg.org/makila/minecraftgo/internal/protocol/packet"
)

type Listener struct{ net.Listener }

// ListenMC listen as TCP but Accept a mc Conn
func ListenMC(addr string) (*Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{l}, nil
}

// Accept a minecraft Conn
func (l Listener) Accept() (Conn, error) {
	conn, err := l.Listener.Accept()
	return Conn{
		Socket: conn,
		Reader: *bufio.NewReader(conn),
		Writer: conn,
	}, err
}

type Conn struct {
	Socket net.Conn
	bufio.Reader
	io.Writer
}

func (c *Conn) Close() {
	c.Socket.Close()
}

func (c *Conn) WritePacket(p packet.Packet) {
	bp, _ := p.Marshal()
	WritePacket(c, int(p.ID()), bp)
}
