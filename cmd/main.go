package main

import (
	"flag"
	"github.com/to77e/password-generator/internal/app/command"
	"log"
)

func main() {
	var (
		create = flag.String("c", "", "create password for service name")
		read   = flag.String("r", "", "read password by service name")
		list   = flag.Bool("l", false, "list of service names")
	)

	flag.Parse()

	if len(*create) > 0 {
		if err := command.CreatePassword(*create); err != nil {
			log.Fatal(err)
		}
	}

	if len(*read) > 0 {
		if err := command.ReadName(*read); err != nil {
			log.Fatal(err)
		}
	}

	if *list {
		if err := command.ListNames(); err != nil {
			log.Fatal(err)
		}
	}
}
