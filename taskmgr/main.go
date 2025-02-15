package main

import (
	"github.com/MrShanks/Taska/taskmgr/server"
	"github.com/MrShanks/Taska/utils"
)

func main() {
	cfg := utils.LoadConfig()

	server.Listen(cfg)
}
