package types

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type UUID struct {
	MostSig  uint64
	LeastSig uint64
}

const (
	SegmentBits = 0x7F
	ContinueBit = 0x80
)

func ReadVarInt(reader *bufio.Reader) (int32, error) {
	var value int32
	var position uint

	for {
		currentByte, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		value |= int32(currentByte&SegmentBits) << position

		if (currentByte & ContinueBit) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, errors.New("VarInt is too big")
		}
	}

	return value, nil
}
func ReadString(reader *bufio.Reader) (string, error) {
	length32, err := ReadVarInt(reader)
	if err != nil {
		return "", err
	}
	length := int(length32)
	fmt.Println(length)

	if length < 0 || length > 32767*4 { // *4 because of UTF-8 byte sizes
		return "", fmt.Errorf("string length %d is out of bounds", length)
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ReadUShort(reader *bufio.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf), nil
}
func ReadArray(r *bufio.Reader) {}

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
