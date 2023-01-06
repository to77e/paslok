package command

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/to77e/password-generator/tools/aes"
	"os"
	"strings"
)

func ReadName(name, cipherKey string) error {
	var (
		tmp        string
		values     []string
		err        error
		file       *os.File
		password   string
		decryptStr string
	)

	if file, err = os.OpenFile(fileName, os.O_RDONLY, 0666); err != nil {
		return fmt.Errorf("failed to open file: %v\n", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tmp = scanner.Text()
		if decryptStr, err = aes.Decrypt(tmp, cipherKey); err != nil {
			return err
		}
		values = strings.Split(decryptStr, ";")
		if values[0] == name {
			password = values[1]
		}
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %v", err)
	}
	if len(password) == 0 {
		return fmt.Errorf("name %v doesn't exist", name)
	}

	if err = clipboard.WriteAll(password); err != nil {
		return fmt.Errorf("faild to clipboard password: %v", err)
	}
	return nil
}
