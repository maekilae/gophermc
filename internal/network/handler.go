package network

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"

	"codeberg.org/makila/minecraftgo/internal/protocol/packets"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Listener struct{ net.Listener }

var (
	bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
)

// ListenMC listen as TCP but Accept a mc Conn
func ListenMC(addr string) (*Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{l}, nil
}

// Accept a minecraft Conn
func (l Listener) Accept() (Handler, error) {
	conn, err := l.Listener.Accept()
	return Handler{
		Socket:       conn,
		Reader:       bufio.NewReader(conn),
		Writer:       bufio.NewWriter(conn),
		isCompressed: false,
		threshold:    32,
	}, err
}

type Handler struct {
	Socket       net.Conn
	Reader       *bufio.Reader
	Writer       *bufio.Writer
	isCompressed bool
	threshold    int32
	State        int
}

func (h *Handler) Close() error {
	if h.Writer != nil {
		h.Writer.Flush()
	}
	return h.Socket.Close()
}
func (h *Handler) WritePacket(p packets.Packet) error {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	tempWriter := bufio.NewWriter(buf)

	// Write the Packet ID as a VarInt
	packetID := types.VarInt(p.ID())
	if err := packetID.Write(tempWriter); err != nil {
		return err
	}

	// Write the Packet Data
	if err := p.Write(tempWriter); err != nil {
		return err
	}
	tempWriter.Flush()
	lengthPrefix := types.VarInt(buf.Len())
	if err := lengthPrefix.Write(h.Writer); err != nil {
		return err
	}

	// Write the payload (ID + Data) to the socket
	if _, err := h.Writer.Write(buf.Bytes()); err != nil {
		return err
	}

	// 7. Flush the TCP stream so the client gets it
	return h.Writer.Flush()
}

// ReadNextPacket reads the framing, identifies the packet, and populates it.
func (h *Handler) ReadNextPacket() (packets.Packet, error) {
	// 1. Read the Total Packet Length
	var length types.VarInt
	if err := length.Read(h.Reader); err != nil {
		return nil, err // This usually means the client disconnected
	}

	// 2. Read the Packet ID
	var packetID types.VarInt
	if err := packetID.Read(h.Reader); err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %w", err)
	}

	// 3. Look up the packet constructor in our registry based on current State
	stateMap, stateExists := ServerBoundRegistry[h.State]
	if !stateExists {
		return nil, fmt.Errorf("unknown connection state: %d", h.State)
	}

	constructor, packetExists := stateMap[int32(packetID)]
	if !packetExists {
		// If we don't know the packet, we must discard the remaining bytes
		// of the payload so it doesn't corrupt the next packet in the stream!
		// payloadLength := int(length) - len(packetID.ToBytes(nil))
		// h.Reader.Discard(payloadLength)
		return nil, fmt.Errorf("unhandled packet ID 0x%02X in state %d", packetID, h.State)
	}

	// 4. Instantiate the empty packet
	packet := constructor()

	// 5. Populate the packet using the interface method!
	if err := packet.Read(h.Reader); err != nil {
		return nil, fmt.Errorf("failed to decode packet 0x%02X: %w", packetID, err)
	}

	return packet, nil
}
