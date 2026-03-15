package encryption

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type Keys struct {
	key *rsa.PrivateKey
	Key rsa.PublicKey
}
type cfb8 struct {
	b         cipher.Block
	blockSize int
	iv        []byte
	decrypt   bool
}

func GenerateRSA() (Keys, error) {

	k, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return Keys{}, err
	}

	return Keys{key: k, Key: k.PublicKey}, nil
}

func (k *Keys) PubKeyToBytes() ([]byte, error) {
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

func NewCFB8Stream(block cipher.Block, iv []byte, decrypt bool) cipher.Stream {
	ivCopy := make([]byte, len(iv))
	copy(ivCopy, iv)
	return &cfb8{
		b:         block,
		blockSize: block.BlockSize(),
		iv:        ivCopy,
		decrypt:   decrypt,
	}
}

func (x *cfb8) XORKeyStream(dst, src []byte) {
	for i := range src {
		val := make([]byte, x.blockSize)
		x.b.Encrypt(val, x.iv)

		c := src[i] ^ val[0]

		copy(x.iv, x.iv[1:])
		if x.decrypt {
			x.iv[x.blockSize-1] = src[i]
		} else {
			x.iv[x.blockSize-1] = c
		}
		dst[i] = c
	}
}
