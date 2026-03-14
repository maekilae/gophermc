package network

import (
	"codeberg.org/makila/minecraftgo/internal/protocol/packets"
)

// Define our connection states
const (
	StateHandshaking = 0
	StateStatus      = 1
	StateLogin       = 2
	StatePlay        = 3
)

// ServerBoundRegistry holds the incoming packets we expect from the client
var ServerBoundRegistry = map[int]map[int32]func() packets.Packet{
	StateHandshaking: {
		0x00: func() packets.Packet { return &packets.Handshake{} },
	},
	StateStatus: {
		0x00: func() packets.Packet { return &packets.StatusRequest{} },
		// 0x01: func() packets.Packet { return &packets.StatusPing{} },
	},
	// Login and Play states would go here...
}
