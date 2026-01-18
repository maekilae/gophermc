package network

import (
	"bufio"
	"crypto/rand"
	"log/slog"
	"net"
	"time"

	"codeberg.org/makila/minecraftgo/internal/protocol/packet"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	nextState := handshake(r)

	if nextState == 1 {
		_, _ = types.ReadVarInt(r) // Length
		packetID, _ := types.ReadVarInt(r)
		if packetID == 0x00 {
			slog.Info("Status Requested")
			StatusResponse(conn)
		} else {
			slog.Info("Ping Requested")
		}
		return
	}

	for {
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, _ = types.ReadVarInt(r)
		pID, _ := types.ReadVarInt(r)
		if pID == 0x00 && nextState != -1 {
			s.loginReq(r, conn)
			nextState = -1
		}
		if pID == 0x01 && nextState == -1 {
			ss, _ := types.ReadByteArray(r)
			t, _ := types.ReadByteArray(r)
			// NOTE REMOVE LOG MSG
			slog.Info("Encryption Response", "Shared Secret", ss, "Token", t)
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

func (s *Server) loginReq(r *bufio.Reader, conn net.Conn) {
	username, _ := types.ReadString(r)

	_, _ = types.ReadUUID(r)
	slog.Info("New login request", "Username", username)
	k, _ := s.Key.PubKeyToBytes()
	t := make([]byte, 4)
	_, _ = rand.Read(t)
	en := packet.Encryption{
		ServerID:   "",
		PubKey:     k,
		Token:      t,
		ShouldAuth: false,
	}
	resp, _ := en.Marshal()
	WritePacket(conn, int(en.ID()), resp)
}
