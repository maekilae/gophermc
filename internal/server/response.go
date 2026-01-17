package server

import (
	"codeberg.org/makila/minecraftgo/internal/packet"
	"encoding/json"
	"net"
)

func StatusResponse(c net.Conn) {
	statusResponse, _ := json.Marshal(map[string]interface{}{
		"version": map[string]interface{}{
			"name":     "1.21.11",
			"protocol": 774,
		},
		"players": map[string]interface{}{
			"max":    100,
			"online": 5,
			"sample": []map[string]string{
				{"name": "Gopher", "id": "00000000-0000-0000-0000-000000000000"},
			},
		},
		"description": map[string]string{
			"text": "§bHello from §aGo§f Server!",
		},
	})

	// Packet ID 0x00: Status Response
	// String is a VarInt length followed by bytes
	respData := append(packet.WriteVarInt(len(statusResponse)), statusResponse...)
	packet.WritePacket(c, 0x00, respData)
}
