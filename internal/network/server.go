package network

import (
	"bufio"
	"log/slog"
	"net"
	"sync"

	"codeberg.org/makila/minecraftgo/internal/db"
	"codeberg.org/makila/minecraftgo/internal/encryption"
)

type Server struct {
	Name       string
	key        encryption.Keys
	addr       string
	tcpaddr    *net.TCPAddr
	port       string
	listener   *Listener
	Players    map[string]net.Conn
	mu         sync.Mutex
	log        *slog.Logger
	IPCRequest chan IPC
}

type IPC interface {
	Send()
	Recive()
}

// ListenMC listen as TCP but Accept a mc Conn
func InitListener(addr string) (*Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{l}, nil
}

// Accept a minecraft Conn
func (s *Server) Accept() (*Handler, error) {
	conn, err := s.listener.Accept()

	return &Handler{
		Socket:       conn,
		Reader:       bufio.NewReader(conn),
		Writer:       bufio.NewWriter(conn),
		isCompressed: false,
		threshold:    32,
		serverKey:    &s.key,
		isEncrypted:  false,
	}, err
}

func NewServer(name, port string, db *db.McDB) *Server {
	ln, err := InitListener(port)
	if err != nil {
		panic(err)
	}
	k, _ := encryption.GenerateRSA()

	return &Server{
		Name:     name,
		key:      k,
		addr:     "127.0.0.1",
		port:     port,
		listener: ln,
		log:      slog.Default().WithGroup("Server"),
	}

}

func (s *Server) RunServer() {
	slog.Info("Starting Server", slog.String("port", s.port))
	defer s.listener.Close()
	for {
		c, e := s.Accept()
		if e == nil {
			go s.HandleConnection(*c)
		}
	}
}
