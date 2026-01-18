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

func WriteString(s string) []byte {
	// Convert the string to a byte slice
	strBytes := []byte(s)

	// 1. Write the length as a VarInt
	buf := WriteVarInt(int(len(strBytes)))
	buf = append(buf, strBytes...)

	return buf

}

func ReadByteArray(r *bufio.Reader) ([]byte, error) {
	// 1. Read the length of the array
	length, err := ReadVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read array length: %w", err)
	}

	// 2. Safety Check: Prevent "Memory Exhaustion" attacks.
	// A malicious client could send a VarInt of 2,147,483,647 to crash your RAM.
	if length < 0 || length > 32767 {
		return nil, fmt.Errorf("prefixed array length out of bounds: %d", length)
	}

	// 3. Read the actual bytes
	data := make([]byte, length)
	_, err = io.ReadFull(r, data) // Use ReadFull to ensure we get exactly 'length' bytes
	if err != nil {
		return nil, fmt.Errorf("failed to read array data: %w", err)
	}

	return data, nil
}

func WriteByteArray(data []byte) []byte {
	buf := WriteVarInt(len(data))
	buf = append(buf, data...)
	return buf
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
