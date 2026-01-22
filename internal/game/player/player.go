package player

import (
	"github.com/google/uuid"
)

type Properties struct {
	Name      string
	Value     string
	Signature string
}

type Player struct {
	Username string
	Uuid     uuid.UUID
	Props    Properties
}
