package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/to77e/password-generator/internal/app/command"
	"log"
	"os"
)

func main() {
	var (
		create    = flag.String("c", "", "create password for service name")
		read      = flag.String("r", "", "read password by service name")
		list      = flag.Bool("l", false, "list of service names")
		cipherKey string
	)

	flag.Parse()
	_ = godotenv.Load()

	cipherKey, found := os.LookupEnv("CIPHERKEY")
	if !found {
		log.Fatalf("failed to load key")
	}

	if len(*create) > 0 {
		if err := command.CreatePassword(cipherKey, *create, flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	}

	if len(*read) > 0 {
		if err := command.ReadName(*read, cipherKey); err != nil {
			log.Fatal(err)
		}
	}

	if *list {
		if err := command.ListNames(cipherKey); err != nil {
			log.Fatal(err)
		}
	}
}
