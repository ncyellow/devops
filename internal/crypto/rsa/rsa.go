package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func parseRSAPrivateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func parseRSAPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("key type is not rsa.PublicKey")
	}
}

func ReadRSAPublicKey(rsaPrivateKeyLocation string) (*rsa.PublicKey, error) {
	public, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		return nil, err
	}
	return parseRSAPublicKeyFromPEM(string(public))
}

func ReadRSAPrivateKey(rsaPrivateKeyLocation string) (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		return nil, err
	}
	return parseRSAPrivateKeyFromPEM(string(priv))
}
