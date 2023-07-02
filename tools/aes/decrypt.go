package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// Decrypt -
func Decrypt(encryptedString, keyString string) (string, error) {
	if encryptedString == "" {
		return "", nil
	}

	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", fmt.Errorf("decode key: %w", err)
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", fmt.Errorf("decode encrypted string: %w", err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, enc := enc[:nonceSize], enc[nonceSize:]

	res, err := gcm.Open(nil, nonce, enc, nil)
	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}

	return string(res), nil
}
