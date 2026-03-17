package handler

import (
	"bufio"
	"net"

	"github.com/maekilae/gophermc/internal/encryption"
	"github.com/maekilae/gophermc/internal/game/player"
	"github.com/maekilae/gophermc/internal/server/ipc"
)

type Handler struct {
	Socket       net.Conn
	Reader       *bufio.Reader
	Writer       *bufio.Writer
	isCompressed bool
	threshold    int32
	State        int
	requestChan  chan ipc.IPC
	replyChan    chan ipc.IPC
	serverKey    *encryption.Keys
	sharedSecret []byte
	verifyToken  []byte
	isEncrypted  bool
	Player       player.PlayerData
}

func (h *Handler) Close() error {
	if h.Writer != nil {
		h.Writer.Flush()
	}
	return h.Socket.Close()
}

// func (h *Handler) HandleConnection() {
// 	defer h.Close()

// 	// nextState := handshake(&h.Reader)
// 	for {
// 		switch h.State {
// 		case StateHandshake:
// 			if h.Handshake() != nil {
// 				return
// 			}
// 		case StateStatus:
// 			if h.Status() != nil {
// 				return
// 			}
// 		case StateLogin:
// 			if h.Login() != nil {
// 				return
// 			}
// 		}
// 	}
// }
