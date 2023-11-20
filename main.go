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
	// case "test_split":
	// 	core.SplitFile("test.pdf", ".", 1024*1024)
	// case "test_recomb":
	// 	core.CombineFiles(".", "./op")
	case "ping":
		core.Client_ping()
	case "test_upload":
		core.Create_File("test.pdf", "./test")
	case "worker":
		core.GetDataNode()
		log.Println("Starting Worker")
	default:
		panic("invalid node type")
	}
}
