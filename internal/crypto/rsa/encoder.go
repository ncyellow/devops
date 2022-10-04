package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type Encoder struct {
	publicKey *rsa.PublicKey
	label     []byte
}

func NewEncoder(publicFileName string) (*Encoder, error) {
	publicKey, err := ReadRSAPublicKey(publicFileName)
	if err != nil {
		return nil, err
	}
	return &Encoder{
		publicKey: publicKey,
		label:     []byte(""),
	}, nil
}

func (e Encoder) Encode(message []byte) ([]byte, error) {
	hash := sha256.New()
	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, e.publicKey, message, e.label)

	if err != nil {
		return nil, err
	}
	return cipherText, nil
}
