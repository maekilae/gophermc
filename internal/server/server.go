package server

import (
	"bufio"
	"log/slog"
	"net"
	"sync"

	"github.com/maekilae/gophermc/config"
	db "github.com/maekilae/gophermc/internal/database"
	"github.com/maekilae/gophermc/internal/encryption"
	"github.com/maekilae/gophermc/internal/server/handler"
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

type Listener struct{ net.Listener }

// ListenMC listen as TCP but Accept a mc Conn
func InitListener(addr string) (*Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{l}, nil
}

// Accept a minecraft Conn
func (s *Server) Accept() (*handler.Handler, error) {
	conn, err := s.listener.Accept()

	return &handler.Handler{
		Socket:      conn,
		Reader:      bufio.NewReader(conn),
		Writer:      bufio.NewWriter(conn),
		Threshold:   32,
		RequestChan: s.IPCRequest,
		ReplyChan:   make(chan ipc.IPC),
		ServerKey:   &s.key,
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
	for req := range s.IPCRequest {
		switch req := req.(type) {

		case *ipc.ServerStatus:
			slog.Info("Server Status Requested")
			// Populate the struct
			req.Name = s.Name
			req.Protocol = int(s.Version.Protocol)
			req.MaxPlayers = s.Properties.MaxPlayers
			req.OnlinePlayers = len(s.Players)
			req.Description = s.Properties.Motd
			req.Favicon = ""
			req.EnforcesSecureChat = s.Properties.EnforceSecureProfile

			// Send it back exclusively to the goroutine that asked for it
			req.Reply <- req

		case *ipc.Blacklist:
			slog.Info("Blacklist Requested")
			// Do logic...
			// req.Reply <- req
		}
	}
}

func (s *Server) RunServer() {
	slog.Info("Starting Server", slog.String("port", s.port))
	s.IPCRequest = make(chan ipc.IPC, 40)
	defer s.listener.Close()
	go s.ipcHandler()
	for {
		h, e := s.Accept()
		if e == nil {
			go h.HandleConnection()
		}
	}
}
