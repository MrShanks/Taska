package utils

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

// ReadVersionFromConfig is a backup way for the program to read the
// version from a config.yaml file.
// The version is passed and backed in at building time but if you
// run the program with go run main.go this function takes care of
// reading the version we are currently in
func ReadVersionFromConfig() string {
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
