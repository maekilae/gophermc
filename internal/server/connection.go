package server

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/maekilae/gophermc/internal/protocol/packets"
	"github.com/maekilae/gophermc/internal/protocol/types"
	"github.com/maekilae/gophermc/internal/server/ipc"
	"github.com/maekilae/gophermc/internal/server/login"
)

func (s *Server) HandleConnection(h Handler) {
	defer h.Close()

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
			serverStatus := ipc.ServerStatus{Reply: h.replyChan}
			serverStatus.Fetch(h.requestChan)
			statusData := packets.StatusResponse{
				Version: packets.Version{
					Name:     serverStatus.Name,
					Protocol: serverStatus.Protocol,
				},
				Players: packets.Players{
					Max:    serverStatus.MaxPlayers,
					Online: serverStatus.OnlinePlayers,
				},
				Description: packets.Description{
					Text: serverStatus.Description,
				},
				EnforcesSecureChat: serverStatus.EnforcesSecureChat,
			}
			h.WritePacket(statusData)

		case *packets.Ping:
			slog.Debug("Ping recived")
			h.WritePacket(pk)
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

			if h.threshold > -1 {
				cPk := packets.Compression{
					Threshold: types.VarInt(h.threshold),
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
