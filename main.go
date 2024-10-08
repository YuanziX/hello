package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/yuanzix/hello/gen"
	"github.com/yuanzix/hello/server"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	pool := &server.Server{
		Chatrooms: make(map[string][]*server.Connection),
	}

	pb.RegisterChatRoomServer(grpcServer, pool)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	fmt.Println("Server started at port :8080")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
}
