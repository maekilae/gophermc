package network

import (
	"log/slog"

	"codeberg.org/makila/minecraftgo/internal/protocol/packets"
)

func (s *Server) HandleConnection(h Handler) {
	defer h.Close()

	// nextState := handshake(&h.Reader)
	pk, err := h.ReadNextPacket()
	if err != nil {
		return
	}
	switch pk := pk.(type) {
	case *packets.Handshake:
		h.State = int(pk.NextState)
	case *packets.StatusRequest:
		if h.State != 1 {
			slog.Error("StatusRequest received in non-Status state", "Packet", pk)
			return
		}
		statusData := packets.StatusResponse{
			Version: packets.Version{
				Name:     "1.20.4",
				Protocol: 765,
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

	default:
		slog.Error("Unknown packet in Handshake state", "Packet", pk)
		return
	}
}

func (s *Server) StatusResponse(h *Handler) {
	for {
		pk, err := h.ReadNextPacket()
		if err != nil {
			return
		}
		switch pk := pk.(type) {
		case *packets.StatusRequest:
			statusData := packets.StatusResponse{
				Version: packets.Version{
					Name:     "1.20.4",
					Protocol: 765,
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

		default:
			slog.Error("Unknown packet in Status state", "Packet", pk)
			return
		}
	}
}

// func handshake(r *bufio.Reader) int {
// 	packetLen, _ := types.ReadVarInt(r) // Total Length
// 	packetID, _ := types.ReadVarInt(r)
// 	if packetID == 0x00 {
// 		protocolVer, _ := types.ReadVarInt(r) // Protocol Version
// 		host, _ := types.ReadString(r)
// 		port, _ := types.ReadUShort(r)
// 		nextState, _ := types.ReadVarInt(r)

// 		slog.Info("Packet", "Len", packetLen, "Protocol.V", protocolVer, "Host", host, "Port", port)
// 		return int(nextState)
// 	}
// 	return -1
// }
