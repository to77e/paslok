package command

import (
	"fmt"
	"github.com/to77e/password-generator/internal/app/generator"
	"github.com/to77e/password-generator/tools/aes"
	"os"
	"strings"
	"time"
)

const (
	fileName = ".pwd"
	length   = 18
)

func CreatePassword(cipherKey, name, comment string) error {
	var (
		err        error
		password   string
		encryptStr string
		file       *os.File
	)

	if file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
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

	if password, err = generator.CreatePassword(length); err != nil {
		return fmt.Errorf("failed to create password: %v\n", err)
	}

	tmp := fmt.Sprintf("%s;%s;%s;%s", name, password, comment, time.Now().Format(time.RFC3339))
	if encryptStr, err = aes.Encrypt(tmp, cipherKey); err != nil {
		return err
	}

	if _, err = file.WriteString(encryptStr + "\n"); err != nil {
		return fmt.Errorf("failed to write in file: %v\n", err)
	}
	return nil
}