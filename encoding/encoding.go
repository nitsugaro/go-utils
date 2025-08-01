package encoding

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func EncodeBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func DecodeBase64URL(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func EncodeHex(data []byte) string {
	return hex.EncodeToString(data)
}

func DecodeHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func EncodeBase32(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

func DecodeBase32(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func EncodeBase32NoPadding(data []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
}

func DecodeBase32NoPadding(s string) ([]byte, error) {
	return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
}
