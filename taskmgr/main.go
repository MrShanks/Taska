package main

import (
	"github.com/MrShanks/Taska/taskmgr/server"
	"github.com/MrShanks/Taska/utils"
)

func main() {
	cfg := utils.LoadConfig("config.yaml")

	server.Listen(cfg)
}
