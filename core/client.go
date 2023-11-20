package core

import (
	"bufio"
	context "context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client_ping() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewPingClient(conn)

	req := PingReq{Ping: "Ping!"}
	res, err := client.Pingpong(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling Ping: %v", err)
	}

	log.Printf("Server response: %v", res.Pong)
}

func SplitFile(inputPath, outputDir string, partSize int64) (n int, err error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	numParts := (fileInfo.Size() + partSize - 1) / partSize

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return 0, err
	}

	for i := int64(0); i < numParts; i++ {
		partFileName := fmt.Sprintf("part_%d", i+1)
		partFilePath := filepath.Join(outputDir, partFileName)

		partFile, err := os.Create(partFilePath)
		if err != nil {
			return 0, err
		}
		defer partFile.Close()

		writer := bufio.NewWriterSize(partFile, 4096)

		partOffset := i * partSize
		partSizeRemaining := partSize
		if partOffset+partSizeRemaining > fileInfo.Size() {
			partSizeRemaining = fileInfo.Size() - partOffset
		}

		_, err = file.Seek(partOffset, io.SeekStart)
		if err != nil {
			return 0, err
		}

		_, err = io.CopyN(writer, file, partSizeRemaining)
		if err != nil {
			return 0, err
		}

		writer.Flush()
	}

	return int(numParts), nil
}

func CombineFiles(inputDir, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	files, err := filepath.Glob(filepath.Join(inputDir, "part_*"))
	if err != nil {
		return err
	}

	for _, partFilePath := range files {
		partFile, err := os.Open(partFilePath)
		if err != nil {
			return err
		}
		defer partFile.Close()

		_, err = io.Copy(outputFile, partFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func Download_File(ip_path string, op_path string) error {
	err := CombineFiles(ip_path, op_path)
	if err != nil {
		return err
	}
	log.Println("Recombined")
	return nil
}

func Create_File(ip_path string, op_path string) {
	partSize := int64(1024 * 1024) // 1 mb per part

	numParts, err := SplitFile(ip_path, op_path, partSize)
	if err != nil {
		fmt.Println("Error splitting file:", err)
		return
	}

	fmt.Println("File split into parts successfully.")

	conn, err := grpc.Dial("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewFileServiceClient(conn)

	for i := 0; i < numParts; i++ {
		partFileName := fmt.Sprintf("part_%d", i+1)
		partFilePath := filepath.Join(op_path, partFileName)

		content, _ := os.ReadFile(partFilePath)

		req := PartialFileRequest{Filename: partFileName, Content: content}
		res, err := client.SendPartialFile(context.Background(), &req)
		if err != nil {
			log.Printf("Error while sending file: %v", err)
		}

		log.Printf("Server response: %v", res.Message)
	}

}
