package network

import (
	"log/slog"
	"net"
	"sync"

	"codeberg.org/makila/minecraftgo/internal/encryption"
)

type Server struct {
	Name string
	Key  encryption.Keys

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
	k, _ := encryption.GenerateRSA()

	return &Server{
		Name: name,
		Key:  k,

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
