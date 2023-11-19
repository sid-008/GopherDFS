package main

import (
	"log"
	"os"

	core "github.com/sid-008/big-data-project/core"
)

func main() {
	nodeType := os.Args[1]
	switch nodeType {
	case "master":
		log.Println("Starting Master")
		core.GetMaster().Start()
	case "ping":
		core.Client_ping()
	case "worker":
		log.Println("Starting Master")
	default:
		panic("invalid node type")
	}
}
