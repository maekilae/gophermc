package types

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func (t *Boolean) Read(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	*t = b != 0
	return nil
}

func (t Boolean) Write(writer *bufio.Writer) error {
	if t {
		return writer.WriteByte(0x01)
	}
	return writer.WriteByte(0x00)
}

// --- Bytes ---
func (t *SignedByte) Read(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	*t = SignedByte(b)
	return err
}

func (t SignedByte) Write(writer *bufio.Writer) error {
	return writer.WriteByte(byte(t))
}

func (t *UnsignedByte) Read(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	*t = UnsignedByte(b)
	return err
}

func (t UnsignedByte) Write(writer *bufio.Writer) error {
	return writer.WriteByte(byte(t))
}

// --- Shorts (16-bit) ---
func (t *SignedShort) Read(reader *bufio.Reader) error {
	var buf [2]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return err
	}
	*t = SignedShort(binary.BigEndian.Uint16(buf[:]))
	return nil
}

func (t SignedShort) Write(writer *bufio.Writer) error {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(t))
	_, err := writer.Write(buf[:])
	return err
}

func (t *UnsignedShort) Read(reader *bufio.Reader) error {
	var buf [2]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return err
	}
	*t = UnsignedShort(binary.BigEndian.Uint16(buf[:]))
	return nil
}

func (t UnsignedShort) Write(writer *bufio.Writer) error {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(t))
	_, err := writer.Write(buf[:])
	return err
}

// --- Ints & Longs (Fixed length) ---
func (t *SignedInt) Read(reader *bufio.Reader) error {
	var buf [4]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return err
	}
	*t = SignedInt(binary.BigEndian.Uint32(buf[:]))
	return nil
}

func (t SignedInt) Write(writer *bufio.Writer) error {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(t))
	_, err := writer.Write(buf[:])
	return err
}

func (t *SignedLong) Read(reader *bufio.Reader) error {
	var buf [8]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return err
	}
	*t = SignedLong(binary.BigEndian.Uint64(buf[:]))
	return nil
}

func (t SignedLong) Write(writer *bufio.Writer) error {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(t))
	_, err := writer.Write(buf[:])
	return err
}

// --- VarInt & VarLong ---
func (t *VarInt) Read(reader *bufio.Reader) error {
	var value int32
	var position uint
	for {
		currentByte, err := reader.ReadByte()
		if err != nil {
			return err
		}
		value |= int32(currentByte&SegmentBits) << position
		if (currentByte & ContinueBit) == 0 {
			*t = VarInt(value)
			return nil
		}
		position += 7
		if position >= 35 {
			return errors.New("VarInt is too big")
		}
	}
}

func (t VarInt) Write(writer *bufio.Writer) error {
	ut := uint32(t)
	for {
		b := byte(ut & 0x7f)
		ut >>= 7
		if ut != 0 {
			if err := writer.WriteByte(b | 0x80); err != nil {
				return err
			}
		} else {
			return writer.WriteByte(b)
		}
	}
}

func (t *VarLong) Read(reader *bufio.Reader) error {
	var value int64
	var position uint
	for {
		currentByte, err := reader.ReadByte()
		if err != nil {
			return err
		}
		value |= int64(currentByte&SegmentBits) << position
		if (currentByte & ContinueBit) == 0 {
			*t = VarLong(value)
			return nil
		}
		position += 7
		if position >= 70 { // Max 10 bytes for 64-bit
			return errors.New("VarLong is too big")
		}
	}
}

func (t VarLong) Write(writer *bufio.Writer) error {
	ut := uint64(t)
	for {
		b := byte(ut & 0x7f)
		ut >>= 7
		if ut != 0 {
			if err := writer.WriteByte(b | 0x80); err != nil {
				return err
			}
		} else {
			return writer.WriteByte(b)
		}
	}
}

// --- Strings & Namespaces ---
func (t *StringN) Read(reader *bufio.Reader) error {
	var length VarInt
	if err := length.Read(reader); err != nil {
		return err
	}

	if length < 0 || length > 32767*4 {
		return fmt.Errorf("string length %d is out of bounds", length)
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return err
	}

	*t = StringN(data)
	return nil
}

func (t StringN) Write(writer *bufio.Writer) error {
	data := []byte(t)
	if err := VarInt(len(data)).Write(writer); err != nil {
		return err
	}
	_, err := writer.Write(data)
	return err
}

func (t *Namespace) Read(reader *bufio.Reader) error {
	var str StringN
	err := str.Read(reader)
	*t = Namespace(str)
	return err
}

func (t Namespace) Write(writer *bufio.Writer) error {
	return StringN(t).Write(writer)
}

// --- Position (Sent as a 64-bit integer) ---
func (t *PackedPosition) Read(reader *bufio.Reader) error {
	var val SignedLong
	err := val.Read(reader)
	*t = PackedPosition(val)
	return err
}

func (t PackedPosition) Write(writer *bufio.Writer) error {
	return SignedLong(t).Write(writer)
}

// --- Angle (1 byte representing 1/256 of a full turn) ---
func (t *Angle) Read(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	*t = Angle(b)
	return err
}

func (t Angle) Write(writer *bufio.Writer) error {
	return writer.WriteByte(byte(t))
}

// --- UUID (16 bytes) ---
func (t *UUID) Read(reader *bufio.Reader) error {
	var buf [16]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return err
	}
	*t = buf
	return nil
}

func (t UUID) Write(writer *bufio.Writer) error {
	_, err := writer.Write(t[:])
	return err
}

// --- ByteArray & BitSets ---
func (t *ByteArray) Read(reader *bufio.Reader) error {
	var length VarInt
	if err := length.Read(reader); err != nil {
		return err
	}

	if length < 0 {
		return errors.New("negative byte array length")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return err
	}
	*t = data
	return nil
}

func (t ByteArray) Write(writer *bufio.Writer) error {
	if err := VarInt(len(t)).Write(writer); err != nil {
		return err
	}
	_, err := writer.Write(t)
	return err
}

func (t *FixedBitSet) Read(reader *bufio.Reader) error {
	var ba ByteArray
	err := ba.Read(reader)
	*t = FixedBitSet(ba)
	return err
}

func (t FixedBitSet) Write(writer *bufio.Writer) error {
	return ByteArray(t).Write(writer)
}
