package login

import (
	"log/slog"

	"github.com/maekilae/gophermc/internal/protocol/packets"
)

type Session interface {
	WritePacket(p packets.Packet) error
	// You can add more later: Disconnect(reason string), GetUsername() string, etc.
}

type HandlerFunc func(session Session, p packets.Packet) error

var handlers = make(map[int32]HandlerFunc)

func Register(id int32, h HandlerFunc) {
	handlers[id] = h
}

func Route(session Session, p packets.Packet) error {
	handler, exists := handlers[p.ID()]
	if !exists {
		slog.Debug("Unhandled login packet", "id", p.ID())
		return nil
	}
	return handler(session, p)
}
