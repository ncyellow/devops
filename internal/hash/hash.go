package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type EncodeFunc func(msg string) string

// CreateEncodeFunc возвращает функцию, принимает текст и возвращает подписанный ключом хеш
func CreateEncodeFunc(secretKey string) EncodeFunc {
	return func(msg string) string {
		if secretKey == "" {
			return ""
		}

		h := hmac.New(sha256.New, []byte(secretKey))
		h.Write([]byte(msg))
		sign := h.Sum(nil)
		return hex.EncodeToString(sign)
	}
}

// CheckSign проверка подписи на корректность
func CheckSign(secretKey string, msg string, correctResult string) bool {
	if secretKey == "" {
		return true
	}
	return msg == correctResult
}
