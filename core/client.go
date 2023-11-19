package core

import (
	context "context"
	"log"

	"google.golang.org/grpc"
)

func Client_ping() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewPingClient(conn)

	req := PingReq{Ping: "Ping!"}
	res, err := client.Pingpong(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Printf("Server response: %v", res.Pong)
}
