package handler

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/maekilae/gophermc/internal/protocol/packets"
	"github.com/maekilae/gophermc/internal/protocol/types"
)

var (
	bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
)

func (h *Handler) WritePacket(p packets.Packet) error {
	payloadBuf := bufPool.Get().(*bytes.Buffer)
	payloadBuf.Reset()
	defer bufPool.Put(payloadBuf)

	tmpWriter := bufio.NewWriter(payloadBuf)
	packetID := types.VarInt(p.ID())
	if err := packetID.Write(tmpWriter); err != nil {
		return err
	}
	if err := p.Write(tmpWriter); err != nil {
		return err
	}
	tmpWriter.Flush()

	finalBuf := bufPool.Get().(*bytes.Buffer)
	finalBuf.Reset()
	defer bufPool.Put(finalBuf)

	finalWriter := bufio.NewWriter(finalBuf)

	if h.isCompressed {
		uncompressedSize := payloadBuf.Len()

		if int32(uncompressedSize) >= h.Threshold {
			dataLength := types.VarInt(uncompressedSize)
			if err := dataLength.Write(finalWriter); err != nil {
				return err
			}
			finalWriter.Flush()

			zWriter := zlib.NewWriter(finalBuf)
			zWriter.Write(payloadBuf.Bytes())
			zWriter.Close()
		} else {
			dataLength := types.VarInt(0)
			if err := dataLength.Write(finalWriter); err != nil {
				return err
			}
			finalWriter.Flush()
			finalBuf.Write(payloadBuf.Bytes())
		}
	} else {
		finalBuf.Write(payloadBuf.Bytes())
	}

	totalLength := types.VarInt(finalBuf.Len())
	if err := totalLength.Write(h.Writer); err != nil {
		return err
	}
	if _, err := h.Writer.Write(finalBuf.Bytes()); err != nil {
		return err
	}

	return h.Writer.Flush()
}

func (h *Handler) ReadNextPacket() (packets.Packet, error) {
	var totalLength types.VarInt
	if err := totalLength.Read(h.Reader); err != nil {
		return nil, err
	}

	packetBytes := make([]byte, totalLength)
	if _, err := io.ReadFull(h.Reader, packetBytes); err != nil {
		return nil, fmt.Errorf("failed to read full packet: %w", err)
	}
	dataReader := bufio.NewReader(bytes.NewReader(packetBytes))

	if h.isCompressed {
		var dataLength types.VarInt
		if err := dataLength.Read(dataReader); err != nil {
			return nil, err
		}

		if dataLength != 0 {
			zReader, err := zlib.NewReader(dataReader)
			if err != nil {
				return nil, fmt.Errorf("zlib error: %w", err)
			}
			defer zReader.Close()

			dataReader = bufio.NewReader(zReader)
		}
	}

	var packetID types.VarInt
	if err := packetID.Read(dataReader); err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %w", err)
	}
	stateMap, stateExists := ServerBoundRegistry[h.State]
	if !stateExists {
		return nil, fmt.Errorf("unknown connection state: %d", h.State)
	}
	constructor, packetExists := stateMap[int32(packetID)]
	if !packetExists {
		return nil, fmt.Errorf("unhandled packet ID 0x%02X in state %d", packetID, h.State)
	}
	packet := constructor()

	if err := packet.Read(dataReader); err != nil {
		return nil, fmt.Errorf("failed to decode packet 0x%02X: %w", packetID, err)
	}

	return packet, nil
}

type cipherReader struct {
	conn   net.Conn
	stream cipher.Stream
}

func (c *cipherReader) Read(p []byte) (n int, err error) {
	n, err = c.conn.Read(p)
	if n > 0 {
		c.stream.XORKeyStream(p[:n], p[:n])
	}
	return n, err
}

type cipherWriter struct {
	conn   net.Conn
	stream cipher.Stream
}

func (c *cipherWriter) Write(p []byte) (n int, err error) {
	enc := make([]byte, len(p))
	c.stream.XORKeyStream(enc, p)
	return c.conn.Write(enc)
}
