package cryptohelper

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
)

type Decrypt func(string) (string, error)

func NewAESDecryptor(key []byte) Decrypt {
	return func(text string) (string, error) {
		decoded, _ := hex.DecodeString(text)
		c, err := aes.NewCipher(key)
		if err != nil {
			return "", err
		}
		res := make([]byte, len(decoded))
		c.Decrypt(res, decoded)

		return string(res), nil
	}
}

func HexGenerator(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}
