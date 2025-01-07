package cryptor

import (
	"fmt"

	"github.com/to77e/paslok/internal/aes"
)

type Service struct {
	cipherKey string
}

func New(passkey string) *Service {
	return &Service{cipherKey: passkey}
}

func (s *Service) Encrypt(password string) (string, error) {
	encryptedPassword, err := aes.Encrypt(password, s.cipherKey)
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}
	return encryptedPassword, nil
}

func (s *Service) Decrypt(encryptedPassword string) (string, error) {
	password, err := aes.Decrypt(encryptedPassword, s.cipherKey)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return password, nil
}
