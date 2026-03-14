package protocol

type Packet interface {
	Unmarshal()
}
