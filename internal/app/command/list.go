package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ListNames() error {
	var (
		tmp    string
		values []string
		err    error
		file   *os.File
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
		values = strings.Split(tmp, ";")
		fmt.Printf("%v\n", values[0])
	}
	return nil
}
