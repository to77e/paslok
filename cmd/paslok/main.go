package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/atotto/clipboard"
	_ "github.com/mattn/go-sqlite3"
	"github.com/to77e/paslok/internal/config"
	"github.com/to77e/paslok/internal/database"
	"github.com/to77e/paslok/internal/models"
	"github.com/to77e/paslok/internal/service/cryptor"
	"github.com/to77e/paslok/internal/service/locker"
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
		slog.With("error", err).Error("init configuration")
		return
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		slog.With("error", err).Error("connect to database")
		return
	}
	defer db.Close() //nolint:errcheck

	lockerService := locker.New(db, cryptor.New(cfg.CipherKey))

	if len(*create) > 0 {
		err = lockerService.Create(*create, flag.Arg(0))
		if err != nil {
			slog.With("error", err).Error("create password")
			return
		}
		fmt.Printf("password for %s created\n", *create)
	}

	if len(*read) > 0 {
		var password string
		password, err = lockerService.Read(*read)
		if err != nil {
			slog.With("error", err).Error("read password")
			return
		}

		if err = clipboard.WriteAll(password); err != nil {
			slog.With("error", err).Error("clipboard password")
			return
		}
	}

	if *list {
		var resources []models.Resource
		resources, err = lockerService.List()
		if err != nil {
			slog.With("error", err).Error("list of names")
			return
		}

		for _, resource := range resources {
			fmt.Printf("%s %s\n", resource.Name, resource.Comment)
		}
	}

	if len(*update) > 0 {
		if err = lockerService.Update(*update, flag.Arg(0), flag.Arg(1)); err != nil {
			slog.With("error", err).Error("update password")
			return
		}
		fmt.Printf("password for %s updated\n", *update)
	}

	if len(*remove) > 0 {
		if err = lockerService.Delete(*remove); err != nil {
			slog.With("error", err).Error("delete password")
			return
		}
		fmt.Printf("password for %s deleted\n", *remove)
	}

	if *ver {
		fmt.Printf("version: %s\n", version)
	}
}
