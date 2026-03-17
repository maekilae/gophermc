package handler

import (
	"bufio"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/maekilae/gophermc/internal/encryption"
	"github.com/maekilae/gophermc/internal/game/player"
	"github.com/maekilae/gophermc/internal/protocol/packets"
	"github.com/maekilae/gophermc/internal/protocol/types"
	"github.com/maekilae/gophermc/internal/server/ipc"
	"github.com/maekilae/gophermc/internal/server/login"
	"github.com/maekilae/gophermc/internal/server/status"
)

type Handler struct {
	Socket       net.Conn
	Reader       *bufio.Reader
	Writer       *bufio.Writer
	isCompressed bool
	Threshold    int32
	State        int
	RequestChan  chan ipc.IPC
	ReplyChan    chan ipc.IPC
	ServerKey    *encryption.Keys
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

func (h *Handler) HandleConnection() {
	defer h.Close()
	h.isEncrypted = false
	h.isCompressed = false

	// nextState := handshake(&h.Reader)
	for {
		switch h.State {
		case StateHandshake:
			if h.Handshake() != nil {
				return
			}
		case StateStatus:
			if h.Status() != nil {
				return
			}
		case StateLogin:
			if h.Login() != nil {
				return
			}
		}
	}
}

func (h *Handler) Handshake() error {
	for {
		pk, err := h.ReadNextPacket()
		if err != nil {
			return err
		}
		switch pk := pk.(type) {
		case *packets.Handshake:
			slog.Debug("Handshake recived")
			h.State = int(pk.NextState)
			return nil
		}
	}
}

func (h *Handler) Status() error {
	for {
		pk, err := h.ReadNextPacket()
		if err != nil {
			return err
		}
		switch pk := pk.(type) {
		case *packets.StatusRequest:
			status.Route(h, pk)

		case *packets.Ping:
			status.Route(h, pk)

		default:
			slog.Error("Unknown packet in Status state", "Packet", pk)
			return errors.New("Unknown packet in Status state")
		}
	}
}

func (h *Handler) Login() error {
	// login.Register(33, &h)
	for {
		pk, err := h.ReadNextPacket()
		if err != nil {
			return err
		}
		switch pk := pk.(type) {
		case *packets.LoginStart:
			err := login.Route(h, pk)
			if err != nil {
				return err
			}
		case *packets.EncryptionResponse:
			err := login.Route(h, pk)
			if err != nil {
				return err
			}

			if h.Threshold > -1 {
				cPk := packets.Compression{
					Threshold: types.VarInt(h.Threshold),
				}
				h.WritePacket(&cPk)
				h.isCompressed = true
			}
			pid, err := uuid.Parse(h.Player.UUID)
			gpp := []types.GameProfileProperty{
				{
					Name:      types.StringN(h.Player.Properties[0].Name),
					Value:     types.StringN(h.Player.Properties[0].Value),
					Signature: types.StringN(h.Player.Properties[0].Signature),
				},
			}
			gp := types.GameProfile{
				Username:   "macke01fcb",
				UUID:       types.UUID(pid),
				Properties: gpp,
			}
			ls := packets.LoginSuccess{GameProfile: gp}
			h.WritePacket(&ls)
		case *packets.LoginAcknowledge:
			slog.Debug("Login acknowledged")
			return nil
		}
	}
}

func (h *Handler) EnableEncryption(sharedSecret []byte) error {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return err
	}

	encStream := encryption.NewCFB8Stream(block, sharedSecret, false)
	decStream := encryption.NewCFB8Stream(block, sharedSecret, true)

	if h.Writer != nil {
		h.Writer.Flush()
	}

	cReader := &cipherReader{conn: h.Socket, stream: decStream}
	cWriter := &cipherWriter{conn: h.Socket, stream: encStream}

	h.Reader = bufio.NewReader(cReader)
	h.Writer = bufio.NewWriter(cWriter)

	h.sharedSecret = sharedSecret
	h.isEncrypted = true

	return nil
}

func (h *Handler) Disconnect(reason string) error {
	slog.Info("User disconnected", "reason", reason)
	return errors.New("User disconnected")
}

func (h *Handler) GetUsername() string {
	return string(h.Player.Username)
}

func (h *Handler) GetServerKey() *encryption.Keys {
	return h.ServerKey
}

func (h *Handler) GetToken() ([]byte, error) {
	if h.verifyToken == nil {
		vt := make([]byte, 16)
		_, err := rand.Read(vt)
		if err != nil {
			return nil, err
		}
		h.verifyToken = vt
	}
	return h.verifyToken, nil
}

func (h *Handler) GetPlayer() player.PlayerData {
	return h.Player
}

func (h *Handler) SetSharedSecret(ss []byte) {
	h.sharedSecret = ss
}

func (h *Handler) UpdatePlayer(p player.PlayerData) {
	h.Player = p
}

func (h *Handler) ServerRequest(req ipc.IPC) error {
	req.ReplyChan(h.ReplyChan)
	err := req.Fetch(h.RequestChan)
	return err
}
