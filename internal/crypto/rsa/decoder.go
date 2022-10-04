package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type Decoder struct {
	privateKey *rsa.PrivateKey
	label      []byte
}

func NewDecoder(privateFileName string) (*Decoder, error) {
	privateKey, err := ReadRSAPrivateKey(privateFileName)
	if err != nil {
		return nil, err
	}
	return &Decoder{
		privateKey: privateKey,
		label:      []byte(""),
	}, nil
}

func (e Decoder) Decode(message []byte) ([]byte, error) {
	hash := sha256.New()
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, e.privateKey, message, e.label)

	if err != nil {
		return nil, err
	}
	return plainText, err
}
