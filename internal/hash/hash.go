package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type EncodeFunc func(msg string) string

func CreateEncodeFunc(secretKey string) EncodeFunc {
	return func(msg string) string {
		if secretKey == "" {
			return ""
		}

		h := hmac.New(sha256.New, []byte(secretKey))
		h.Write([]byte(msg))
		sign := h.Sum(nil)
		fmt.Println(hex.EncodeToString(sign))
		return hex.EncodeToString(sign)
	}
}

func CheckSign(secretKey string, msg string, correctResult string) bool {
	if secretKey == "" {
		return true
	}

	return msg == correctResult
}