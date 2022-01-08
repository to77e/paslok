package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/to77e/go-password-generator/internal/app/generator"
)

const (
	length = 18
)

const (
	layout = "2006-01-02"
)

func main() {

	password, err := generator.CreatePassword(length)
	if err != nil {
		log.Fatalf("failed to create password: %v\n", err)
	}

	file, err := os.OpenFile("pwds", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open file: %v\n", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	commentPtr := flag.String("c", "", "comment")
	flag.Parse()
	_, err = file.WriteString(strings.TrimSpace(fmt.Sprintf("%s %s %s", time.Now().Format(layout), password, *commentPtr)) + "\n")
	if err != nil {
		log.Fatalf("failed to write in file: %v\n", err)
	}
}
