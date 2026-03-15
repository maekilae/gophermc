package network

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"sync"

	"codeberg.org/makila/minecraftgo/internal/encryption"
	"codeberg.org/makila/minecraftgo/internal/game/player"
	"codeberg.org/makila/minecraftgo/internal/protocol/packets"
	"codeberg.org/makila/minecraftgo/internal/protocol/types"
)

type Listener struct{ net.Listener }

var (
	bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
)

type Handler struct {
	Socket       net.Conn
	Reader       *bufio.Reader
	Writer       *bufio.Writer
	isCompressed bool
	threshold    int32
	State        int
	IPC          chan IPC
	serverKey    *encryption.Keys
	sharedSecret []byte
	verifyToken  []byte
	isEncrypted  bool
	Player       player.Player
}

func (h *Handler) Close() error {
	if h.Writer != nil {
		h.Writer.Flush()
	}
	return h.Socket.Close()
}
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

		if int32(uncompressedSize) >= h.threshold {
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

func (h *Handler) EnableEncryption(sharedSecret []byte) error {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return err
	}

	encStream := encryption.NewCFB8Stream(block, sharedSecret, false)
	decStream := encryption.NewCFB8Stream(block, sharedSecret, true)

	if h.Writer != nil {
		h.Writer.Flush()
	}

	cReader := &cipherReader{conn: h.Socket, stream: decStream}
	cWriter := &cipherWriter{conn: h.Socket, stream: encStream}

	h.Reader = bufio.NewReader(cReader)
	h.Writer = bufio.NewWriter(cWriter)

	h.sharedSecret = sharedSecret
	h.isEncrypted = true

	return nil
}
