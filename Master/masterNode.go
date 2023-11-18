package master

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Block struct {
	ID        int
	Locations []string
}

type Metadata_entry struct {
	Name              string
	ReplicationFactor uint32
	Id                string
	Size              int64
	Blocks            []Block
}

type MasterNode struct {
	api  *gin.Engine
	meta map[string]*Metadata_entry // {k,v} = {path, data}
}

func (n *MasterNode) Init() (err error) {
	n.api = gin.Default()
	log.Println("Contents of meta", n.meta)

	n.api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})

	return nil
}

func (n *MasterNode) Start() {
	_ = n.api.Run(":9001")
}

var masterNode *MasterNode

func GetMaster() *MasterNode {
	if masterNode == nil {
		masterNode = &MasterNode{}

		if err := masterNode.Init(); err != nil {
			log.Fatal(err)
		}
	}
	return masterNode
}
