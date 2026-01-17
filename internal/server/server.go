package server

import (
	"bytes"
	"codeberg.org/makila/minecraftgo/internal/client"
	"codeberg.org/makila/minecraftgo/internal/packet"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {

	//Basic server name for readability
	Name string

	//Network Address for our server
	addr    string
	tcpaddr *net.TCPAddr

	//Our listener
	listenerConfig net.ListenConfig
	listener       net.Listener

	nrClients uint16
	clients   []client.Client

	//Concurrency, Context & sync
	srvWG     *sync.WaitGroup
	srvCtx    context.Context
	srvCancel context.CancelFunc
	mu        sync.Mutex
}

func NewServer(name, protocol string, port string) Server {
	ln, err := net.Listen(protocol, port)
	if err != nil {
		panic(err)
	}

	return Server{
		Name: name,

		addr: "127.0.0.1",

		listener:  ln,
		nrClients: 0,
	}

}

func (s *Server) RunServer() {}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// 1. READ HANDSHAKE (Packet ID 0x00)
	_, _ = packet.ReadVarInt(conn) // Total Length
	packetID, _ := packet.ReadVarInt(conn)
	if packetID == 0x00 {
		packet.ReadVarInt(conn) // Protocol Version
		hostLen, _ := packet.ReadVarInt(conn)
		host := make([]byte, hostLen)
		io.ReadFull(conn, host) // Server Address
		var port uint16
		binary.Read(conn, binary.BigEndian, &port)
		nextState, _ := packet.ReadVarInt(conn)

		if nextState != 1 {
			loginReq(conn)
		} else {

			_, _ = packet.ReadVarInt(conn) // Length
			if packetID == 0x00 {
				fmt.Printf("[%s] Status Requested\n", conn.RemoteAddr())
				StatusResponse(conn)
			}
			return
		}
	}
	if packetID == 0x01 {
		var payload int64
		binary.Read(conn, binary.BigEndian, &payload)
		fmt.Printf("[%s] Ping Received: %d\n", conn.RemoteAddr(), payload)
		buf := new(bytes.Buffer)

		// 1. Packet ID (0x01)
		// We write it as a VarInt (which for 0x01 is just 0x01)
		buf.WriteByte(0x01)

		// 2. Payload (Long)
		// The spec says: Signed 64-bit integer, two's complement.
		// binary.Write handles int64 (Long) correctly in BigEndian.
		binary.Write(buf, binary.BigEndian, payload)

		// 3. Construct final packet with Length Header
		// Total length = ID (1 byte) + Long (8 bytes) = 9 bytes
		finalPacket := make([]byte, 0)
		finalPacket = append(finalPacket, 0x09) // Length VarInt
		finalPacket = append(finalPacket, buf.Bytes()...)

		conn.Write(finalPacket)
	}
}

func loginReq(r io.Reader) {
	nameLen, _ := packet.ReadVarInt(r)
	nameBuf := make([]byte, nameLen)
	io.ReadFull(r, nameBuf)
	username := string(nameBuf)

	u, _ := packet.ReadUUID(r)
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(u.MostSig>>32),     // First 8 hex chars
		uint16(u.MostSig>>16),     // Next 4
		uint16(u.MostSig),         // Next 4
		uint16(u.LeastSig>>48),    // Next 4
		u.LeastSig&0xFFFFFFFFFFFF) // Final 12
	fmt.Printf(`Username: %s`, username)
	fmt.Printf("UUID: %s", uuid)

}
