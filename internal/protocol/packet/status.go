package packet

import "encoding/json"

type StatusResp struct {
	Json string
}

func (s *StatusResp) ID() int32 {
	return 0x00
}

func (s *StatusResp) Marshal() ([]byte, error) {
	return json.Marshal(s.Json)
}
