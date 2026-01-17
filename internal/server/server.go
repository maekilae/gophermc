package server

import (
	"bufio"
	"log/slog"
	"net"
	"sync"
	"time"

	"codeberg.org/makila/minecraftgo/internal/packet"
)

type Server struct {
	Name string

	addr    string
	tcpaddr *net.TCPAddr
	port    string

	listener net.Listener

	players map[string]net.Conn

	mu sync.Mutex

	log *slog.Logger
}

func NewServer(name, protocol string, port string) *Server {
	ln, err := net.Listen(protocol, port)
	if err != nil {
		panic(err)
	}

	return &Server{
		Name: name,

		addr: "127.0.0.1",
		port: port,

		listener: ln,

		log: slog.Default().WithGroup("Server"),
	}

}

func (s *Server) RunServer() {
	slog.Info("Starting Server", slog.String("port", s.port))
	defer s.listener.Close()
	for {
		c, e := s.listener.Accept()
		if e == nil {
			go s.HandleConnection(c)
		}
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	nextState := handshake(conn, r)

	if nextState == 1 {
		_, _ = packet.ReadVarInt(r) // Length
		packetID, _ := packet.ReadVarInt(r)
		if packetID == 0x00 {
			slog.Info("Status Requested")
			StatusResponse(conn)
		}
		if packetID == 0x01 {
			slog.Info("Ping Requested")
		}
		return
	}

	for {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, _ = packet.ReadVarInt(r)
		pID, _ := packet.ReadVarInt(r)
		if pID == 0x00 && nextState != -1 {
			loginReq(r)
			nextState = -1
		}
	}
}

func handshake(conn net.Conn, r *bufio.Reader) int {
	packetLen, _ := packet.ReadVarInt(r) // Total Length
	packetID, _ := packet.ReadVarInt(r)
	if packetID == 0x00 {
		protocolVer, _ := packet.ReadVarInt(r) // Protocol Version
		host, _ := packet.ReadString(r)
		port, _ := packet.ReadUShort(r)
		nextState, _ := packet.ReadVarInt(r)

		slog.Info("Packet", "Len", packetLen, "Protocol.V", protocolVer, "Host", host, "Port", port)
		return int(nextState)
	}
	return -1
}

func loginReq(r *bufio.Reader) {
	username, _ := packet.ReadString(r)

	_, _ = packet.ReadUUID(r)
	slog.Info("New login request", "Uname", username)

}
