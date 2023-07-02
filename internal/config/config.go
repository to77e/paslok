package config

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var config embed.FS

type Config struct {
	FilePath  string `yaml:"filePath"`
	CipherKey string `yaml:"cipherKey"`
}

func (c *Config) ReadConfig() error {
	yamlFile, err := config.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("read config.yaml: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return fmt.Errorf("unmarshal config.yaml: %w", err)
	}
	return nil
}
