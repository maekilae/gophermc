package server

import (
	"bufio"
	"log/slog"
	"net"
	"sync"

	"github.com/maekilae/gophermc/config"
	db "github.com/maekilae/gophermc/internal/database"
	"github.com/maekilae/gophermc/internal/encryption"
	"github.com/maekilae/gophermc/internal/server/ipc"
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
	IPCRequest chan ipc.IPC
	Version    config.ServerVersion
	Properties config.ServerProperties
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

func NewServer(name, port string, db *db.McDB, version config.ServerVersion, properties config.ServerProperties) *Server {
	ln, err := InitListener(port)
	if err != nil {
		panic(err)
	}
	k, _ := encryption.GenerateRSA()

	return &Server{
		Name:       name,
		key:        k,
		addr:       "127.0.0.1",
		port:       port,
		listener:   ln,
		log:        slog.Default().WithGroup("Server"),
		Version:    version,
		Properties: properties,
	}
}

func (s *Server) ipcHandler() {
	for {
		select {
		case req := <-s.IPCRequest:
			switch req := req.(type) {
			case ipc.ServerStatus:
				slog.Info("Server Status Requested")
				req.Name = s.Name
				req.Protocol = int(s.Version.Protocol)
				req.MaxPlayers = s.Properties.MaxPlayers
				req.OnlinePlayers = len(s.Players)
				req.Description = s.Properties.Motd
				req.Favicon = ""
				req.EnforcesSecureChat = s.Properties.EnforceSecureProfile
				s.IPCRequest <- req
			}
		}
	}
}

func (s *Server) RunServer() {
	slog.Info("Starting Server", slog.String("port", s.port))
	s.IPCRequest = make(chan ipc.IPC, 40)
	defer s.listener.Close()
	for {
		c, e := s.Accept()
		if e == nil {
			go s.HandleConnection(*c)
		}
	}
}
