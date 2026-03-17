package handler

import "github.com/maekilae/gophermc/internal/protocol/packets"

const (
	StateHandshake = 0
	StateStatus    = 1
	StateLogin     = 2
	StatePlay      = 3
)

// ServerBoundRegistry holds the incoming packets we expect from the client
var ServerBoundRegistry = map[int]map[int32]func() packets.Packet{
	StateHandshake: {
		0x00: func() packets.Packet { return &packets.Handshake{} },
	},
	StateStatus: {
		0x00: func() packets.Packet { return &packets.StatusRequest{} },
		0x01: func() packets.Packet { return &packets.Ping{} },
	},
	StateLogin: {
		0x00: func() packets.Packet { return &packets.LoginStart{} },
		0x01: func() packets.Packet { return &packets.EncryptionResponse{} },
		0x03: func() packets.Packet { return &packets.LoginAcknowledge{} },
	},
	// Login and Play states would go here...
}
