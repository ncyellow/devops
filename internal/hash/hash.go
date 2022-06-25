package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type EncodeFunc func(msg string) string

func CreateEncodeFunc(secretKey string) EncodeFunc {
	return func(msg string) string {
		if secretKey == "" {
			return ""
		}

		data := hex.EncodeToString([]byte(msg))

		h := hmac.New(sha256.New, []byte(secretKey))
		h.Write([]byte(data))
		sign := h.Sum(nil)
		return hex.EncodeToString(sign)
	}
}

func CheckSign(secretKey string, msg string, correctResult string) bool {
	if secretKey == "" {
		return true
	}

	return msg == correctResult
}
