package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/to77e/password-generator/internal/app"
	"log"
	"os"
	"strings"
	"time"
)

const (
	length = 18
)

const (
	fileName = "pwd.csv"
)

var (
	cNamePtr = flag.String("c", "", "create password with name")
	rNamePtr = flag.String("r", "", "read password by name")
	listPtr  = flag.Bool("l", false, "list of names")
)

// TODO: crypt output file

func main() {
	var err error
	flag.Parse()

	if len(*cNamePtr) > 0 {
		if err = writeStringToFile(); err != nil {
			log.Fatal(err)
		}
	}

	if *listPtr {
		if err = listNamesFromFile(); err != nil {
			log.Fatal(err)
		}
	}

	if len(*rNamePtr) > 0 {
		if err = readNameFromFile(); err != nil {
			log.Fatal(err)
		}
	}
}

func writeStringToFile() error {
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

	if password, err = app.CreatePassword(length); err != nil {
		return fmt.Errorf("failed to create password: %v\n", err)
	}

	newStr := fmt.Sprintf("%s;%s;%s\n", *cNamePtr, password, time.Now().Format(time.RFC3339))
	if _, err = file.WriteString(newStr); err != nil {
		return fmt.Errorf("failed to write in file: %v\n", err)
	}
	return nil
}

func listNamesFromFile() error {
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

func readNameFromFile() error {
	var (
		tmp      string
		values   []string
		err      error
		file     *os.File
		password string
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
		if values[0] == *rNamePtr {
			password = values[1]
		}
	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %v", err)
	}
	if err = clipboard.WriteAll(password); err != nil {
		return fmt.Errorf("faild to clipboard password: %v", err)
	}
	return nil
}
