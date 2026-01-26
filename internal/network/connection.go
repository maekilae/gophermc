package network

import (
	"bufio"
	"crypto/md5"
	"log/slog"
	"time"

	"codeberg.org/makila/minecraftgo/internal/protocol/types"

	"github.com/google/uuid"
)

func (s *Server) HandleConnection(conn Conn) {
	defer conn.Close()

	nextState := handshake(&conn.Reader)

	if nextState == 1 {
		_, _ = types.ReadVarInt(&conn.Reader) // Length
		packetID, _ := types.ReadVarInt(&conn.Reader)
		if packetID == 0x00 {
			slog.Info("Status Requested")
			StatusResponse(&conn)
		} else {
			slog.Info("Ping Requested")
		}
		return
	}

	for {
		conn.Socket.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, _ = types.ReadVarInt(&conn.Reader)
		pID, _ := types.ReadVarInt(&conn.Reader)
		if pID == 0x00 && nextState != -1 {
			_, _, _ = s.StartLogin(&conn)
			nextState = -1
		}
	}
}

func handshake(r *bufio.Reader) int {
	packetLen, _ := types.ReadVarInt(r) // Total Length
	packetID, _ := types.ReadVarInt(r)
	if packetID == 0x00 {
		protocolVer, _ := types.ReadVarInt(r) // Protocol Version
		host, _ := types.ReadString(r)
		port, _ := types.ReadUShort(r)
		nextState, _ := types.ReadVarInt(r)

		slog.Info("Packet", "Len", packetLen, "Protocol.V", protocolVer, "Host", host, "Port", port)
		return int(nextState)
	}
	return -1
}


func NameToUUID(name string) uuid.UUID {
	version := 3
	h := md5.New()
	h.Write([]byte("OfflinePlayer:"))
	h.Write([]byte(name))
	var id uuid.UUID
	h.Sum(id[:0])
	id[6] = (id[6] & 0x0f) | uint8((version&0xf)<<4)
	id[8] = (id[8] & 0x3f) | 0x80 // RFC 4122 variant
	return id
}
