package login

import (
	"bytes"

	"github.com/maekilae/gophermc/internal/auth"
	"github.com/maekilae/gophermc/internal/protocol/packets"
)

func init() {
	Register(0x01, handleEncryptionResponse)
}
func handleEncryptionResponse(session Session, p packets.Packet) error {
	encResp, ok := p.(*packets.EncryptionResponse)
	if !ok {
		return session.Disconnect("Invalid packet")
	}
	sk := session.GetServerKey()
	ss, err := sk.Decrypt(encResp.SharedSecret)
	if err != nil {
		return session.Disconnect("Failed to decrypt shared secret")
	}
	vt, err := sk.Decrypt(encResp.VerifyToken)
	if err != nil {
		return session.Disconnect("Failed to decrypt verify token")
	}
	cvt, err := session.GetToken()
	if err != nil {
		return session.Disconnect("Failed to client token(Server side)")
	}
	if !bytes.Equal(cvt, vt) {
		return session.Disconnect("Token mismatch")
	}
	k, err := sk.PubKeyToBytes()
	if err != nil {
		return session.Disconnect("Failed to marshal public key")
	}
	player := session.GetPlayer()

	pData, err := auth.SendHash(string(player.Username), auth.AuthDigest("", ss, k))
	if err != nil {
		return session.Disconnect("Failed to send hash to Mojang")
	}

	player.Properties = pData.Properties
	session.UpdatePlayer(player)

	if err := session.EnableEncryption(ss); err != nil {
		return session.Disconnect("Failed to enable encryption")
	}

	return nil
}
