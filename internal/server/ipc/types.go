package ipc

import "time"

type IPC interface{}

type ServerStatus struct {
	Name          string
	Protocol      int
	MaxPlayers    int
	OnlinePlayers int
	// Players            []Player
	Description        string
	Favicon            string
	EnforcesSecureChat bool
}

type Blacklist struct {
	Players map[string]time.Duration
}
