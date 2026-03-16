package login

import (
	"encoding/hex"
	"log/slog"

	"github.com/maekilae/gophermc/internal/game/player"
	"github.com/maekilae/gophermc/internal/protocol/packets"
)

func init() {
	// Register the LoginStart packet ID (usually 0x00 in the Login state)
	Register(0x00, handleLoginStart)
}
func handleLoginStart(session Session, p packets.Packet) error {
	loginStart, ok := p.(*packets.LoginStart)
	if !ok {
		return session.Disconnect("Invalid packet")
	}
	slog.Info("Player logging in", "name", loginStart.Username)

	sk := session.GetServerKey()
	k, err := sk.PubKeyToBytes()
	if err != nil {
		return session.Disconnect("Failed to get public key")
	}
	token, err := session.GetToken()
	if err != nil {
		return session.Disconnect("Failed to get token")
	}
	encReq := &packets.EncryptionRequest{
		ServerID:   "",
		PubKey:     k,
		Token:      token,
		ShouldAuth: true,
	}
	if err := session.WritePacket(encReq); err != nil {
		return session.Disconnect("Failed to send encryption request")
	}
	player := player.PlayerData{
		Username:   string(loginStart.Username),
		UUID:       hex.EncodeToString(loginStart.UUID[:]),
		Properties: []player.Property{},
	}
	session.UpdatePlayer(player)
	return nil
}
