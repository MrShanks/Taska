package main

import (
	"log"
	"os"

	"github.com/MrShanks/Taska/taskmgr/server"
	"github.com/MrShanks/Taska/utils"
)

func main() {
	cfg := utils.LoadConfig("config.yaml")

	err := server.Listen(cfg)
	if err != nil {
		log.Printf("Error starting server: %v", err)
		os.Exit(1)
	}
}
