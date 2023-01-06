package command

import (
	"fmt"
	"github.com/to77e/password-generator/internal/app/generator"
	"os"
	"time"
)

const (
	fileName = "pwd.csv"
	length   = 18
)

func CreatePassword(name string) error {
	var (
		err      error
		password string
		file     *os.File
	)

	if file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return fmt.Errorf("failed to open file: %v\n", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	if password, err = generator.CreatePassword(length); err != nil {
		return fmt.Errorf("failed to create password: %v\n", err)
	}

	newStr := fmt.Sprintf("%s;%s;%s\n", name, password, time.Now().Format(time.RFC3339))
	if _, err = file.WriteString(newStr); err != nil {
		return fmt.Errorf("failed to write in file: %v\n", err)
	}
	return nil
}
