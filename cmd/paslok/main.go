package main

import (
	"flag"
	"log"

	"github.com/to77e/paslok/internal/app/command"
	"github.com/to77e/paslok/internal/config"
)

func main() {
	var (
		create = flag.String("c", "", "create password for service name")
		read   = flag.String("r", "", "read password by service name")
		list   = flag.Bool("l", false, "list of service names")
		cfg    config.Config
	)
	flag.Parse()

	err := cfg.ReadConfig(&cfg)
	if err != nil {
		log.Fatal("init configuration: %w", err)
	}

	if len(*create) > 0 {
		err = command.CreatePassword(cfg.CipherKey, cfg.FilePath, *create, flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(*read) > 0 {
		if err = command.ReadName(*read, cfg.CipherKey, cfg.FilePath); err != nil {
			log.Fatal(err)
		}
	}

	if *list {
		err = command.ListNames(cfg.CipherKey, cfg.FilePath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
