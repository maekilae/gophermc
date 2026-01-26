package network

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"

	packet "codeberg.org/makila/minecraftgo/internal/protocol/clientbound"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Listener struct{ net.Listener }

var (
	bufPool  = sync.Pool{New: func() any { return new(bytes.Buffer) }}
)

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
		Socket:       conn,
		Reader:       *bufio.NewReader(conn),
		Writer:       conn,
		isCompressed: false,
		threshold:    32,
	}, err
}

type Conn struct {
	Socket net.Conn
	bufio.Reader
	io.Writer
	isCompressed bool
	threshold    int32
}

func (c *Conn) Close() {
	c.Socket.Close()
}

func createHeader(id int, data []byte) (buf []byte) {
	_ = types.VarInt(1+len(data)).ToBytes(buf)
	fmt.Println(buf)
	_= types.VarInt(id).ToBytes(buf)
	return
}

func (c *Conn) WritePacket(p packet.Packet) {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	var data []byte
	if c.isCompressed == true {
		data, _ = p.Marshal()

	} else {
		data, _ = p.Marshal()

	}
	buf.Write(createHeader(int(p.ID()),data))
	buf.Write(data)
	c.Write(buf.Bytes())
}

func (c *Conn) ReadPacket() (int32, int32) {
	size, _ := types.ReadVarInt(&c.Reader)
	id, _ := types.ReadVarInt(&c.Reader)
	return id, size

}
