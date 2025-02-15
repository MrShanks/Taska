package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Spec        struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"spec"`
}

func LoadConfig() *Config {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Failed to open the config file: %v", err)
		os.Exit(1)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Printf("Failed to parse the config file: %v", err)
		os.Exit(1)
	}

	return cfg
}
