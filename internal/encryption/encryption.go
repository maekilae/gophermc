package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type Keys struct {
	key *rsa.PrivateKey
	Key rsa.PublicKey
}

func GenerateRSA() (Keys, error) {

	k, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return Keys{}, err
	}

	return Keys{key: k, Key: k.PublicKey}, nil
}

func (k *Keys) PubKeyToBytes() ([]byte, error) {
	// MarshalPKIXPublicKey returns the DER-encoded public key
	pubBytes, err := x509.MarshalPKIXPublicKey(&k.Key)
	if err != nil {
		return nil, err
	}
	return pubBytes, nil
}
func (k *Keys) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, k.key, data)
}

func ParseKeyStream(data []byte) (any, error) {
	return x509.ParsePKIXPublicKey(data)
}
