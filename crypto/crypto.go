package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"hash"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.NewString()
}

func GetRandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func HashSHA1(input string) []byte {
	return hashBytes(sha1.New(), input)
}

func HashSHA256(input string) []byte {
	return hashBytes(sha256.New(), input)
}

func HashSHA384(input string) []byte {
	return hashBytes(sha512.New384(), input)
}

func HashSHA512(input string) []byte {
	return hashBytes(sha512.New(), input)
}

func hashBytes(h hash.Hash, input string) []byte {
	h.Write([]byte(input))
	return h.Sum(nil)
}

func HmacSHA1(input string, base64urlKey string) ([]byte, error) {
	return hmacBytes(sha1.New, input, base64urlKey)
}

func HmacSHA256(input string, base64urlKey string) ([]byte, error) {
	return hmacBytes(sha256.New, input, base64urlKey)
}

func HmacSHA384(input string, base64urlKey string) ([]byte, error) {
	return hmacBytes(sha512.New384, input, base64urlKey)
}

func HmacSHA512(input string, base64urlKey string) ([]byte, error) {
	return hmacBytes(sha512.New, input, base64urlKey)
}

func hmacBytes(f func() hash.Hash, input string, base64urlKey string) ([]byte, error) {
	key, err := base64.RawURLEncoding.DecodeString(base64urlKey)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(f, key)
	mac.Write([]byte(input))
	return mac.Sum(nil), nil
}
