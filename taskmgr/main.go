package main

import (
	"github.com/MrShanks/Taska/taskmgr/logger"
	"github.com/MrShanks/Taska/taskmgr/server"
)

func main() {
	logger.InitLogger()
	server.Listen()
}
