package packet

import (
	"encoding/binary"
	"errors"
	"io"
)

type UUID struct {
	MostSig  uint64
	LeastSig uint64
}

// ReadVarInt reads a Minecraft-style variable length integer.
func ReadVarInt(r io.Reader) (int, error) {
	var v uint32
	for i := 0; ; i++ {
		var b [1]byte
		if _, err := r.Read(b[:]); err != nil {
			return 0, err
		}
		v |= uint32(b[0]&0x7F) << uint32(7*i)
		if b[0]&0x80 == 0 {
			break
		}
		if i >= 5 {
			return 0, errors.New("VarInt is too big")
		}
	}
	return int(v), nil
}

// WriteVarInt encodes an int into a byte slice.
func WriteVarInt(value int) []byte {
	var res []byte
	uValue := uint32(value)
	for {
		b := byte(uValue & 0x7F)
		uValue >>= 7
		if uValue != 0 {
			b |= 0x80
		}
		res = append(res, b)
		if uValue == 0 {
			break
		}
	}
	return res
}

// ReadUUID reads 16 bytes and splits them into two uint64s
func ReadUUID(r io.Reader) (UUID, error) {
	var u UUID
	err := binary.Read(r, binary.BigEndian, &u.MostSig)
	if err != nil {
		return u, err
	}
	err = binary.Read(r, binary.BigEndian, &u.LeastSig)
	return u, err
}

// WriteUUID writes the two uint64s as a 16-byte block
func WriteUUID(w io.Writer, u UUID) error {
	err := binary.Write(w, binary.BigEndian, u.MostSig)
	if err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, u.LeastSig)
}
