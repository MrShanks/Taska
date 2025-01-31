package cmd

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
	} `yaml:"app"`
}

func readVersionFromConfig() string {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Failed to open the config file: %v", err)
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Printf("Failed to parse the config file: %v", err)
		os.Exit(1)
	}

	return cfg.App.Version
}
