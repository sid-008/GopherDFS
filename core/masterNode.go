package core

import (
	context "context"
	"log"
	"net"

	"google.golang.org/grpc"
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
	srv  *grpc.Server
	lis  net.Listener
	meta map[string]*Metadata_entry // {k,v} = {path, metadata}
}

type Pingserver struct {
	UnimplementedPingServer
}

func (s *Pingserver) Pingpong(ctx context.Context, req *PingReq) (*PongResp, error) {
	log.Println("Recvd. a ping, sent back pong.")
	return &PongResp{Pong: "Pong"}, nil
}

func (n *MasterNode) Init() (err error) {
	log.Println("Contents of meta", n.meta)
	n.lis, err = net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}

	n.srv = grpc.NewServer()
	RegisterPingServer(n.srv, &Pingserver{})

	log.Println("Server listening on :9090")

	if err := n.srv.Serve(n.lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return nil
}

func (n *MasterNode) Start() {
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
