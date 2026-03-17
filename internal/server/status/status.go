package status

import (
	"github.com/maekilae/gophermc/internal/protocol/packets"
	"github.com/maekilae/gophermc/internal/server/ipc"
)

func init() {
	Register(0x00, handleStatusRequest)
	Register(0x01, handlePing)
}
func handleStatusRequest(session Session, p packets.Packet) error {
	serverStatus := ipc.ServerStatus{}
	err := session.ServerRequest(&serverStatus)
	if err != nil {
		return err
	}
	statusData := packets.StatusResponse{
		Version: packets.Version{
			Name:     serverStatus.Name,
			Protocol: serverStatus.Protocol,
		},
		Players: packets.Players{
			Max:    serverStatus.MaxPlayers,
			Online: serverStatus.OnlinePlayers,
		},
		Description: packets.Description{
			Text: serverStatus.Description,
		},
		EnforcesSecureChat: serverStatus.EnforcesSecureChat,
	}
	err = session.WritePacket(statusData)
	if err != nil {
		return err
	}
	return nil

}

func handlePing(session Session, p packets.Packet) error {
	return session.WritePacket(p)
}
