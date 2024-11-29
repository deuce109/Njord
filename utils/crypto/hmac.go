package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func HMACSHA512(message string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))

	hash.Write([]byte(message))

	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func HMACEquals(a string, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}
