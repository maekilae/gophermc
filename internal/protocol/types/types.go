package types

import (
	"bufio"

	"github.com/google/uuid"
)

const (
	SegmentBits = 0x7F
	ContinueBit = 0x80
)

type Types interface {
	Read(*bufio.Reader)
	ToBytes([]byte)
}

// https://minecraft.wiki/w/Java_Edition_protocol/Data_types
type (
	Boolean        bool
	SignedByte     int8
	UnsignedByte   uint8
	SignedShort    int16
	UnsignedShort  uint16
	SignedInt      int32
	SignedLong     int64
	StringN        string
	Namespace      string
	VarInt         int32
	VarLong        int64
	PackedPosition int64 // Marshaled version of Position struct
	Angle          uint8
	UUID           uuid.UUID
	FixedBitSet    ByteArray

	ByteArray []byte

	// Float
	// Double

)
type TextFormat struct{}

type TextEvents struct{}

// https://minecraft.wiki/w/Text_component_format
type TextComponent struct {
	Content     string
	Children    []TextComponent
	Format      TextFormat
	Interaction TextEvents
}

// https://minecraft.wiki/w/Java_Edition_protocol/Entity_metadata#Entity_Metadata_Format
type EntityMetaData struct {
}

type ComponentData struct {
}

// https://minecraft.wiki/w/Java_Edition_protocol/Slot_data
type SlotData struct {
	HasItem                    Boolean
	ItemCount                  VarInt
	ItemID                     VarInt
	NumberOfComponentsToAdd    VarInt
	NumberOfComponentsToRemove VarInt
	// Components
}

type NBT struct{}

type Position struct {
	X int32
	Y int16
	Z int32
}
