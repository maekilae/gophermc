package encryption

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateRSA() (*rsa.PrivateKey, error) {

	k, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	return k, nil
}
