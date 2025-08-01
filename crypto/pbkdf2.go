package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"

	"github.com/nitsugaro/go-utils/encoding"
	"golang.org/x/crypto/pbkdf2"
)

type ALG_SHA string

const (
	SHA256 ALG_SHA = "sha256"
	SHA384 ALG_SHA = "sha384"
	SHA512 ALG_SHA = "sha512"
)

type Pbkdf2 struct {
	Password   string  `json:"password"`
	Iterations int     `json:"iterations"`
	SaltSize   int     `json:"salt_size"`
	KeyLen     int     `json:"key_len"`
	Alg        ALG_SHA `json:"alg"`
}

func IsPbkdf2Hash(hash string) bool {
	return strings.HasPrefix(hash, "pbkdf2_")
}

func Pbkdf2Hash(data *Pbkdf2) (string, error) {
	if data.SaltSize == 0 {
		data.SaltSize = 16
	}

	salt := make([]byte, data.SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hashFunc, err := getHashFunc(data.Alg)
	if err != nil {
		return "", err
	}

	if data.Iterations == 0 {
		data.Iterations = 1000
	}

	if data.KeyLen == 0 {
		data.KeyLen = 32
	}

	hash := pbkdf2.Key([]byte(data.Password), salt, data.Iterations, data.KeyLen, hashFunc)

	// pbkdf2_sha256$<iterations>$<base64_salt>$<base64_hash>
	prefix := fmt.Sprintf("pbkdf2_%s", strings.ToLower(string(data.Alg)))
	result := fmt.Sprintf("%s$%d$%s$%s", prefix, data.Iterations, encoding.EncodeBase64(salt), encoding.EncodeBase64(hash))
	return result, nil
}

func VerifyPbkdf2Hash(password string, hashStored string) bool {
	parts := strings.Split(hashStored, "$")
	if len(parts) != 4 {
		return false
	}

	prefix := parts[0]
	var alg ALG_SHA
	if strings.HasPrefix(prefix, "pbkdf2_") {
		alg = ALG_SHA(strings.TrimPrefix(prefix, "pbkdf2_"))
	} else {
		return false
	}

	hashFunc, err := getHashFunc(alg)
	if err != nil {
		return false
	}

	iterations := 0
	_, err = fmt.Sscanf(parts[1], "%d", &iterations)
	if err != nil {
		return false
	}

	salt, err := encoding.DecodeBase64(parts[2])
	if err != nil {
		return false
	}

	expectedHash, err := encoding.DecodeBase64(parts[3])
	if err != nil {
		return false
	}

	keyLen := len(expectedHash)
	calculated := pbkdf2.Key([]byte(password), salt, iterations, keyLen, hashFunc)

	return hmacEqual(expectedHash, calculated)
}

func getHashFunc(alg ALG_SHA) (func() hash.Hash, error) {
	switch strings.ToLower(string(alg)) {
	case "sha256":
		return sha256.New, nil
	case "sha384":
		return sha512.New384, nil
	case "sha512":
		return sha512.New, nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm: %s", alg)
	}
}

func hmacEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}
