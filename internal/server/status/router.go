package status

import (
	"log/slog"

	"github.com/maekilae/gophermc/internal/protocol/packets"
	"github.com/maekilae/gophermc/internal/server/ipc"
)

type Session interface {
	WritePacket(p packets.Packet) error
	Disconnect(reason string) error
	ServerRequest(req ipc.IPC) error
}

type HandlerFunc func(session Session, p packets.Packet) error

var handlers = make(map[int32]HandlerFunc)

func Register(id int32, h HandlerFunc) {
	handlers[id] = h
}

func Route(session Session, p packets.Packet) error {
	handler, exists := handlers[p.ID()]
	if !exists {
		slog.Debug("Unhandled status packet", "id", p.ID())
		return nil
	}
	return handler(session, p)
}
