package packets

import (
	"bufio"
)

type Packet interface {
	// ID returns the protocol ID for this packet (e.g., 0x00 for Handshake)
	ID() int32

	// Read populates the packet's fields from the reader.
	// It assumes the packet header (length and ID) has already been read.
	Read(r *bufio.Reader) error

	// Write writes the packet's fields to the writer.
	// It assumes the framing will be added later by the Handler.
	Write(w *bufio.Writer) error
}
