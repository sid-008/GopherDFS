package core

import (
	context "context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

type DataNode struct {
	ID string
	// Fs   afero.Fs
	srv  *grpc.Server
	lis  net.Listener
	root string
}

// type Directory struct {
// 	Name    string
// 	Files   []File
// 	SubDirs []Directory
// }
//
// type File struct {
// 	Name     string
// 	Contents []byte
// }

func (dn *DataNode) CreateRoot(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

type FileServiceserver struct {
	UnimplementedFileServiceServer
}

var nn_count int

func (n *DataNode) Init(i int, wg *sync.WaitGroup) {
	nn_count += 1

	n.lis, _ = net.Listen("tcp", fmt.Sprintf(":909%d", i))

	// dataNode := NewDataNode(afero.NewOsFs())
	dataNode := DataNode{}

	// log.Println(nn_count)
	dataNode.CreateRoot(fmt.Sprintf("/tmp/node%d", i))

	n.root = fmt.Sprintf("/tmp/node%d", i)

	n.srv = grpc.NewServer()
	RegisterFileServiceServer(n.srv, &FileServiceserver{})

	log.Printf("Server listening on :909%d", i)

	if err := n.srv.Serve(n.lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

var wg sync.WaitGroup

func GetDataNode() {
	dn := DataNode{}
	wg.Add(1)
	go dn.Init(1, &wg)
	wg.Add(1)
	go dn.Init(2, &wg)
	wg.Wait()
}

func gen_ran(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func (fs *FileServiceserver) SendPartialFile(ctx context.Context, req *PartialFileRequest) (res *PartialFileResponse, err error) {
	log.Println("Worker recvd. ", req.Filename)

	res = &PartialFileResponse{}

	fd, err := os.OpenFile(fmt.Sprintf("/tmp/node%d/%s", gen_ran(1, nn_count), req.Filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return res, err
	}
	_, err = fd.Write(req.Content)
	if err != nil {
		res.Success = false
		res.Message = "Error"
		return res, err
	}
	defer fd.Close()

	_, _ = fd.Write(req.Content)

	res.Success = true
	res.Message = "Done"

	return res, nil
}
