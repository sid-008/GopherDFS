package main

import (
	"log"
	"os"

	master "github.com/sid-008/big-data-project/Master"
)

func main() {
	nodeType := os.Args[1]
	switch nodeType {
	case "master":
		log.Println("Starting Master")
		master.GetMaster().Start()
	case "worker":
		log.Println("Starting Master")
	default:
		panic("invalid node type")
	}
}
