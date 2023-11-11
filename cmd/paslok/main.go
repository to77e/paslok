package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	_ "github.com/mattn/go-sqlite3"
	"github.com/to77e/paslok/internal/config"
	"github.com/to77e/paslok/internal/database"
	"github.com/to77e/paslok/internal/service/cryptor"
	"github.com/to77e/paslok/internal/service/locker"
	"log"
)

var version string

func main() {
	var (
		create = flag.String("create", "", "create password by name")
		read   = flag.String("read", "", "read password by name")
		remove = flag.String("delete", "", "delete password by name")
		update = flag.String("update", "", "update password by name")
		list   = flag.Bool("list", false, "list of names")
		ver    = flag.Bool("version", false, "show version")
		cfg    config.Config
	)

	flag.Parse()

	err := cfg.ReadConfig(&cfg)
	if err != nil {
		log.Fatal("init configuration: %w", err)
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("connect to database: %s", err.Error())
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	lockerService := locker.New(db, cryptor.New(cfg.CipherKey))

	if len(*create) > 0 {
		err = lockerService.Create(*create, flag.Arg(0))
		if err != nil {
			log.Fatalf("create password: %s", err.Error())
		}
		fmt.Printf("password for %s created\n", *create)
	}

	if len(*read) > 0 {
		password, err := lockerService.Read(*read)
		if err != nil {
			log.Fatalf("read password: %s", err.Error())
		}

		if err = clipboard.WriteAll(password); err != nil {
			log.Fatalf("faild to clipboard password: %v", err)
		}
	}

	if *list {
		resources, err := lockerService.List()
		if err != nil {
			log.Fatalf("list of names: %s", err.Error())
		}

		for _, resource := range resources {
			fmt.Printf("%s %s\n", resource.Name, resource.Comment)
		}
	}

	if len(*update) > 0 {
		if err = lockerService.Update(*update, flag.Arg(0), flag.Arg(1)); err != nil {
			log.Fatalf("update password: %s", err.Error())
		}
		fmt.Printf("password for %s updated\n", *update)
	}

	if len(*remove) > 0 {
		if err = lockerService.Delete(*remove); err != nil {
			log.Fatalf("delete password: %s", err.Error())
		}
		fmt.Printf("password for %s deleted\n", *remove)
	}

	if *ver {
		fmt.Printf("version: %s\n", version)
	}
}
