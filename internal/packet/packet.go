package packet

// WritePacket wraps a packet ID and data with its length.
func WritePacket(id int, data []byte) []byte {
	idBytes := WriteVarInt(id)
	lenBytes := WriteVarInt(len(idBytes) + len(data))
	return append(lenBytes, append(idBytes, data...)...)
}
