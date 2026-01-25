package packet

type State int

type Packet interface {
	ID() int32
	Marshal() ([]byte, error)    // Convert struct to bytes
	Unmarshal(data []byte) error // Convert bytes to struct
}
