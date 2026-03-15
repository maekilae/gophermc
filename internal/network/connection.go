package network

import (
	"bytes"
	"crypto/rand"
	"errors"
	"log/slog"

	"codeberg.org/makila/minecraftgo/internal/api"
	"codeberg.org/makila/minecraftgo/internal/protocol/packets"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
	"github.com/google/uuid"
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
			statusData := packets.StatusResponse{
				Version: packets.Version{
					Name:     "1.21.11",
					Protocol: 774,
				},
				Players: packets.Players{
					Max:    100,
					Online: 5,
				},
				Description: packets.Description{
					Text: "Welcome to my custom Go server!",
				},
				EnforcesSecureChat: false,
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
	for {
		pk, err := h.ReadNextPacket()
		if err != nil {
			return err
		}
		switch pk := pk.(type) {
		case *packets.LoginStart:
			slog.Debug("Login start packet recived from", "Username", pk.Username)
			key, err := h.serverKey.PubKeyToBytes()
			if err != nil {
				return err
			}
			vt := make([]byte, 16)
			_, err = rand.Read(vt)
			if err != nil {
				return err
			}
			h.verifyToken = vt
			enPk := packets.EncryptionRequest{
				ServerID:   "",
				PubKey:     key,
				Token:      vt,
				ShouldAuth: true,
			}
			h.WritePacket(&enPk)
		case *packets.EncryptionResponse:
			ss, err := h.serverKey.Decrypt(pk.SharedSecret)
			if err != nil {
				return err
			}
			vt, err := h.serverKey.Decrypt(pk.VerifyToken)
			if err != nil {
				return err
			}
			if !bytes.Equal(h.verifyToken, vt) {
				return errors.New("Could not authenicate user")
			}
			key, err := h.serverKey.PubKeyToBytes()
			if err != nil {
				return err
			}
			pData, err := api.SendHash("macke01fcb", api.AuthDigest("", ss, key))
			if err != nil {
				return err
			}
			slog.Debug("Player with", "UUID", pData.Properties[0].Signature)

			h.EnableEncryption(ss)

			if h.threshold > -1 {
				cPk := packets.Compression{
					Threshold: types.VarInt(h.threshold),
				}
				h.WritePacket(&cPk)
				h.isCompressed = true
			}
			pid, err := uuid.Parse(pData.ID)
			gpp := []types.GameProfileProperty{{Name: types.StringN(pData.Properties[0].Name), Value: types.StringN(pData.Properties[0].Value), Signature: types.StringN(pData.Properties[0].Signature)}}
			gp := types.GameProfile{Username: "macke01fcb", UUID: types.UUID(pid), Properties: gpp}
			ls := packets.LoginSuccess{GameProfile: gp}
			h.WritePacket(&ls)
		case *packets.LoginAcknowledge:
			slog.Debug("Login acknowledged")
			return nil
		}
	}
}
