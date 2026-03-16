package auth

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/maekilae/gophermc/internal/game/player"
)

func SendHash(u string, h string) (player.PlayerData, error) {
	var content player.PlayerData
	params := url.Values{}
	params.Add("serverId", h)
	params.Add("username", u)

	au := "https://sessionserver.mojang.com/session/minecraft/hasJoined?" + params.Encode()
	// fmt.Println("Calling Mojang:", au)

	resp, err := http.Get(au)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()

	// CRITICAL FIX: Mojang returns 204 No Content if the hash doesn't match!
	if resp.StatusCode == http.StatusNoContent {
		return content, fmt.Errorf("mojang rejected login: player %s not verified", u)
	}
	if resp.StatusCode != http.StatusOK {
		return content, fmt.Errorf("mojang api returned unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}

	// CRITICAL FIX: Actually capture and return the JSON unmarshal error
	err = json.Unmarshal(body, &content)
	return content, err
}

func AuthDigest(serverId string, sharedSecret []byte, publicKey []byte) string {
	h := sha1.New()
	h.Write([]byte(serverId))
	h.Write(sharedSecret)
	h.Write(publicKey)
	hash := h.Sum(nil)

	// Check for "negative" hashes (first bit of first byte is 1)
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Minecraft trims leading zeros and uses lowercase hex
	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// twosComplement performs a manual bitwise two's complement on the byte slice.
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i]) // bitwise NOT
		if carry {
			carry = (p[i] == 0xff)
			p[i]++
		}
	}
	return p
}
