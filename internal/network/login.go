package network

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"

	"codeberg.org/makila/minecraftgo/internal/api"
	"codeberg.org/makila/minecraftgo/internal/protocol/packet"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"

	"github.com/google/uuid"
)

type LoginHandler interface {
	StartLogin(conn *Conn) (uname string, id uuid.UUID)
}

func (s *Server) StartLogin(conn *Conn) (uname string, id uuid.UUID, err error) {
	// r := bufio.NewReader(*conn)
	uname, _ = types.ReadString(&conn.Reader)
	_, _ = types.ReadUUID(&conn.Reader)
	slog.Info("New login request", "Username", uname)
	k, _ := s.Key.PubKeyToBytes()
	vt := make([]byte, 16)
	_, _ = rand.Read(vt)
	en := packet.Encryption{
		ServerID:   "",
		PubKey:     k,
		Token:      vt,
		ShouldAuth: true,
	}
	resp, _ := en.Marshal()
	WritePacket(conn, int(en.ID()), resp)
	_, _ = conn.ReadPacket()
	// _, _ = types.ReadVarInt(&conn.Reader)
	// _, _ = types.ReadVarInt(&conn.Reader)
	ss, _ := types.ReadByteArray(&conn.Reader)
	t, _ := types.ReadByteArray(&conn.Reader)
	// NOTE REMOVE LOG MSG
	ss, _ = s.Key.Decrypt(ss)
	t, _ = s.Key.Decrypt(t)
	if !bytes.Equal(t, vt) {
		return "", uuid.UUID{}, errors.New("Client token mismatch")
	}
	slog.Info("Encryption Response", "Shared Secret", ss, "Token", t)
	k, _ = s.Key.PubKeyToBytes()
	pd, e := api.SendHash(uname, api.AuthDigest("", ss, k))
	if e != nil {
		slog.Error("Could not authenticate with mojang")
	}
	fmt.Println(pd)
	WritePacket(conn, 0x02, pd.Marshal())
	return uname, uuid.UUID{}, err
}
