package network

import (

	packet "codeberg.org/makila/minecraftgo/internal/protocol/clientbound"
)

func StatusResponse(c *Conn) {
	status := packet.ServerStatus{
    Version: packet.Version{
        Name:     "1.20.1",
        Protocol: 763,
    },
    Players: packet.Players{
        Max:    100,
        Online: 2,
        Sample: []packet.Player{
            {Name: "Alice", ID: "u-u-i-d-1"},
            {Name: "Bob", ID: "u-u-i-d-2"},
        },
    },
    Description: packet.Description{
        Text: "Welcome to the best Minecraft server!",
    },
    Favicon:            "data:image/png;base64,iVBOR...",
    EnforcesSecureChat: true,
}

	// Packet ID 0x00: Status Response
	// String is a VarInt length followed by bytes
	// respData := append(types.WriteVarInt(len(statusResponse)), statusResponse...)

	c.WritePacket(status)
}
