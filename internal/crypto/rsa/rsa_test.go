package rsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRSAPublicKey(t *testing.T) {
	publicKey, err := ReadRSAPublicKey("test_data/rsa.public")
	assert.NoError(t, err)
	assert.NotNil(t, publicKey)

	publicKey, err = ReadRSAPublicKey("test_data/rsa.private")
	assert.Error(t, err)
	assert.Nil(t, publicKey)

	publicKey, err = ReadRSAPublicKey("test_data/rsa.nofile")
	assert.Error(t, err)
	assert.Nil(t, publicKey)
}

func TestReadRSAPrivateKey(t *testing.T) {
	privateKey, err := ReadRSAPrivateKey("test_data/rsa.private")
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)

	privateKey, err = ReadRSAPrivateKey("test_data/rsa.public")
	assert.Error(t, err)
	assert.Nil(t, privateKey)

	privateKey, err = ReadRSAPrivateKey("test_data/rsa.nofile")
	assert.Error(t, err)
	assert.Nil(t, privateKey)
}

func TestEncoderDecoder(t *testing.T) {
	encoder, err := NewEncoder("test_data/rsa.public")
	assert.NoError(t, err)
	assert.NotNil(t, encoder)

	decoder, err := NewDecoder("test_data/rsa.private")
	assert.NoError(t, err)
	assert.NotNil(t, decoder)

	want := []byte("simple test")

	cipherText, err := encoder.Encode(want)
	assert.NoError(t, err)
	assert.NotNil(t, cipherText)

	decodeMsg, err := decoder.Decode(cipherText)
	assert.NoError(t, err)
	assert.NotNil(t, decodeMsg)

	assert.Equal(t, want, decodeMsg)

}
