package main

import (
	"github.com/gustavoppereira/llv/agent/configuration"
	"github.com/gustavoppereira/llv/agent/manager"
)

// REGEX TO EXTRACT ONLY MESSAGE -> (\d{4}\/\d{2}\/\d{1,2})\s(\d{2}:\d{2}:\d{2})\s(.*)

func main() {
	var agent = manager.NewManager()
	agent.StartAgent(configuration.GetServerToken())
}
