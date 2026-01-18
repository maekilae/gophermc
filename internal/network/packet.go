package network

import (
	"net"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Header struct {
	length    []byte
	packet_id []byte
}

type Packet struct {
	header Header
	data   []byte
}

func createPacket(id int, data []byte) Packet {
	h := Header{
		length:    types.WriteVarInt(1 + len(data)),
		packet_id: types.WriteVarInt(id),
	}
	return Packet{
		header: h,
		data:   data,
	}

}

// WritePacket wraps a packet ID and data with its length.
func WritePacket(c net.Conn, id int, data []byte) {
	p := createPacket(id, data)
	size := len(p.header.length) + len(p.header.packet_id) + len(p.data)
	buffer := make([]byte, 0, size)
	buffer = append(buffer, p.header.length...)
	buffer = append(buffer, p.header.packet_id...)
	buffer = append(buffer, p.data...)
	c.Write(buffer)
}
