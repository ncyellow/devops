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

	publicKey, err = ReadRSAPublicKey("test_data/rsa_with_error.public")
	assert.Error(t, err, "failed to parse PEM block containing the public key")
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

	privateKey, err = ReadRSAPrivateKey("test_data/rsa_with_error.private")
	assert.Error(t, err)
	assert.Nil(t, privateKey, "failed to parse PEM block containing the private key")
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

	// Отдельно тестируем граничные случаи, когда и кодировщику и декодировщику переданы не корректные ключи
	encoder, err = NewEncoder("test_data/rsa.private")
	assert.Error(t, err)
	assert.Nil(t, encoder)

	decoder, err = NewDecoder("test_data/rsa.public")
	assert.Error(t, err)
	assert.Nil(t, decoder)
}
