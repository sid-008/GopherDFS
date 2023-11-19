package core

import "time"

type DataNode struct {
	ID            string    // Unique identifier for the DataNode
	Hostname      string    // Hostname or IP address of the DataNode
	Port          int       // Port number on which the DataNode is running
	StoragePath   string    // Directory path where data blocks are stored
	LastHeartbeat time.Time // Timestamp of the last heartbeat from the DataNode
}
