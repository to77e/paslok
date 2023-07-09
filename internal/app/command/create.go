package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/to77e/paslok/internal/app/generator"
	"github.com/to77e/paslok/tools/aes"
)

func CreatePassword(cipherKey, filePath, name, comment string) error {
	const (
		perm   = 0600
		length = 18
	)

	file, err := os.OpenFile(filepath.Clean(filePath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(perm))
	if err != nil {
		return fmt.Errorf("failed to open file: %v\n", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	if strings.ContainsAny(name, ";") {
		return fmt.Errorf("the name should not contain \";\": %v\n", err)
	}

	password, err := generator.CreatePassword(length)
	if err != nil {
		return fmt.Errorf("failed to create password: %v\n", err)
	}

	tmp := fmt.Sprintf("%s;%s;%s;%s", name, password, comment, time.Now().Format(time.RFC3339))
	encryptStr, err := aes.Encrypt(tmp, cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	if _, err = file.WriteString(encryptStr + "\n"); err != nil {
		return fmt.Errorf("failed to write in file: %v\n", err)
	}
	return nil
}
