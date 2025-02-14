package main

import (
	"github.com/MrShanks/Taska/taskmgr/server"
	"github.com/MrShanks/Taska/utils"
)

var version = "undefined"

func main() {
	if version == "undefined" {
		version = utils.ReadVersionFromConfig()
	}

	server.Listen(version)
}
