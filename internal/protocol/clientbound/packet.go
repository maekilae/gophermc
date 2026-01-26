package clientbound

type State int

type Packet interface {
	ID() int32
	Marshal() ([]byte, error)    // Convert struct to bytes
}
