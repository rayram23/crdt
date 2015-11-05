package main

import (
	"github.com/rayram23/crdt/addonly/addonly"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	addonly.RegisterAddOnlyServer(s, addonly.NewServer())
	s.Serve(lis)
}
