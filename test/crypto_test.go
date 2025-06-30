package test

import (
	"testing"

	"github.com/nitsugaro/go-utils/cipher"
	"github.com/nitsugaro/go-utils/crypto"
	"github.com/nitsugaro/go-utils/encoding"
)

func TestCryptoFunctions(t *testing.T) {
	plainText := "1234"
	hashSha256Base64 := "A6xnQhbz4Vx2HuGl4lXwZ5U2I8iziLRFnhP5eNfIRvQ="

	hashResult := encoding.EncodeBase64(crypto.HashSHA256(plainText))
	if hashResult != hashSha256Base64 {
		t.Errorf("expected hash sha256 encode base64 of '1234' be '%s' and got: %s", hashSha256Base64, hashResult)
	}

	randomKey, _ := crypto.GetRandBytes(32)
	cipherText, err := cipher.EncryptAESGCM(randomKey, []byte(plainText))
	if err != nil {
		t.Errorf("encrypt AES-GCM failure: %s", err.Error())
	}

	plainTextDecrypted, err := cipher.DecryptAESGCM(randomKey, cipherText)
	if err != nil {
		t.Errorf("decrypt AES-GCM failure: %s", err.Error())
	}

	if encoding.EncodeBase64URL(plainTextDecrypted) != encoding.EncodeBase64URL([]byte(plainText)) {
		t.Errorf("wrong decryption AES-GCM")
	}

	cipher1, _ := cipher.EncryptAESCBC(randomKey, []byte(plainText))
	cipher2, _ := cipher.EncryptAESCBC(randomKey, []byte(plainText))

	if encoding.EncodeHex(cipher1) == encoding.EncodeHex(cipher2) {
		t.Errorf("wrong encryption AES-CBC")
	}

	plainText1, _ := cipher.DecryptAESCBC(randomKey, cipher1)
	plainText2, _ := cipher.DecryptAESCBC(randomKey, cipher2)

	if encoding.EncodeHex(plainText1) != encoding.EncodeHex(plainText2) {
		t.Errorf("wrong decryption AES-CBC")
	}
}
