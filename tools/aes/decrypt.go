package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// Decrypt -
func Decrypt(encryptedString, keyString string) (out string, err error) {
	var (
		key, enc []byte
	)

	if encryptedString == "" {
		return "", nil
	}

	if key, err = hex.DecodeString(keyString); err != nil {
		return
	}
	if enc, err = hex.DecodeString(encryptedString); err != nil {
		return
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}

	nonceSize := gcm.NonceSize()
	nonce, enc := enc[:nonceSize], enc[nonceSize:]

	res, err := gcm.Open(nil, nonce, enc, nil)
	if err != nil {
		return
	}

	return fmt.Sprintf("%s", res), nil
}
