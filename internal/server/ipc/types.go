package ipc

import (
	"fmt"
	"time"
)

type IPC interface {
	Fetch(severChan chan<- IPC) error
	ReplyChan(replyChan chan IPC)
}

// ==========================================
// SERVER STATUS
// ==========================================

type ServerStatus struct {
	Reply              chan IPC // Back to the generic IPC channel
	Name               string
	Protocol           int
	MaxPlayers         int
	OnlinePlayers      int
	Description        string
	Favicon            string
	EnforcesSecureChat bool
}

func (s *ServerStatus) Fetch(serverChan chan<- IPC) error {
	// 1. Initialize the generic channel
	s.Reply = make(chan IPC, 1)

	// 2. Send the struct value to the server loop
	serverChan <- s

	// 3. Receive the generic IPC interface
	res := <-s.Reply

	// 4. Type-assert it back to ServerStatus and assign it to itself
	if typedRes, ok := res.(*ServerStatus); ok {
		s = typedRes
		return nil
	}

	return fmt.Errorf("expected ServerStatus, got %T", res)
}

func (s *ServerStatus) ReplyChan(replyChan chan IPC) {
	s.Reply = replyChan
}

// ==========================================
// BLACKLIST
// ==========================================

type Blacklist struct {
	Reply   chan IPC
	Players map[string]time.Duration
}

func (b *Blacklist) Fetch(serverChan chan<- IPC) error {
	b.Reply = make(chan IPC, 1)
	serverChan <- b

	res := <-b.Reply
	if typedRes, ok := res.(*Blacklist); ok {
		b = typedRes
		return nil
	}

	return fmt.Errorf("expected Blacklist, got %T", res)
}
func (b *Blacklist) ReplyChan(replyChan chan IPC) {
	b.Reply = replyChan
}
