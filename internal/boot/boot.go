package boot

import (
	"fmt"

	"github.com/to77e/paslok/internal/config"
	"github.com/to77e/paslok/internal/database"
	"github.com/to77e/paslok/internal/service/cryptor"
	"github.com/to77e/paslok/internal/service/locker"
)

type App struct {
	Config        config.Config
	Database      *database.Database
	LockerService *locker.Service
}

func Initialize() (*App, error) {
	var cfg config.Config

	err := cfg.ReadConfig(&cfg)
	if err != nil {
		return nil, fmt.Errorf("init configuration: %w", err)
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	c := cryptor.New(cfg.CipherKey)
	lockerService := locker.New(db, c)

	return &App{
		Config:        cfg,
		Database:      db,
		LockerService: lockerService,
	}, nil
}
