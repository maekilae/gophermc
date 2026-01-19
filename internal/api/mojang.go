package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PlayerData struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

func SendHash(u string, h string) (PlayerData, error) {
	var content PlayerData
	params := url.Values{}
	params.Add("username", u)
	params.Add("serverId", h)
	au := "https://sessionserver.mojang.com/session/minecraft/hasJoined?" + params.Encode()
	fmt.Println(au)
	resp, err := http.Get(au)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	json.Unmarshal(body, &content)
	return content, err
}

func AuthDigest(serverId string, sharedSecret []byte, publicKey []byte) string {
	h := sha1.New()
	io.WriteString(h, serverId)
	h.Write(sharedSecret)
	h.Write(publicKey)
	hash := h.Sum(nil)

	// Check for "negative" hashes (first bit of first byte is 1)
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Minecraft trims leading zeros and uses lowercase hex
	res := strings.TrimLeft(hex.EncodeToString(hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// twosComplement performs a manual bitwise two's complement on the byte slice.
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = ^p[i] // bitwise NOT
		if carry {
			carry = (p[i] == 0xff)
			p[i]++
		}
	}
	return p
}
