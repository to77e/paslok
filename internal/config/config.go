package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	DBPath    string `env:"PASLOK_DB_PATH" envDefault:"./tmp/paslok.db"`
	CipherKey string `env:"PASLOK_CIPHER_KEY"`
}

func (c *Config) ReadConfig(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	if strings.HasPrefix(cfg.DBPath, "~/") {
		home, _ := os.UserHomeDir()
		cfg.DBPath = filepath.Join(home, cfg.DBPath[2:])
	}

	return nil
}
