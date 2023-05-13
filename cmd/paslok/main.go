package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/to77e/password-generator/internal/app/command"
)

func main() {
	var (
		create    = flag.String("c", "", "create password for service name")
		read      = flag.String("r", "", "read password by service name")
		list      = flag.Bool("l", false, "list of service names")
		cipherKey string
	)

	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading .env file")
	}

	cipherKey, found := os.LookupEnv("CIPHERKEY")
	if !found {
		log.Fatalf("load cipher key")
	}

	filePath, found := os.LookupEnv("FILEPATH")
	if !found {
		log.Fatalf("load file path")
	}

	if len(*create) > 0 {
		if err := command.CreatePassword(cipherKey, filePath, *create, flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	}

	if len(*read) > 0 {
		if err := command.ReadName(*read, cipherKey, filePath); err != nil {
			log.Fatal(err)
		}
	}

	if *list {
		if err := command.ListNames(cipherKey, filePath); err != nil {
			log.Fatal(err)
		}
	}
}
