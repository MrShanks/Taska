package main

import (
	"os"

	"github.com/MrShanks/Taska/taskmgr/server"
	"github.com/MrShanks/Taska/utils"
)

func main() {
	cfg := utils.LoadConfig("config.yaml")

	err := server.Listen(cfg)
	if err != nil {
		os.Exit(1)
	}
}
