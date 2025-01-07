package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func Encrypt(stringToEncrypt, keyString string) (out string, err error) {
	var key []byte

	if stringToEncrypt == "" {
		return "", nil
	}

	if key, err = hex.DecodeString(keyString); err != nil {
		return
	}
	enc := []byte(stringToEncrypt)

	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	res := gcm.Seal(nonce, nonce, enc, nil)
	return fmt.Sprintf("%x", res), nil
}
