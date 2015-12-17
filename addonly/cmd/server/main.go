package main

import (
	"github.com/ianschenck/envflag"
	"github.com/rayram23/crdt/addonly/server"
	"github.com/rayram23/crdt/addonly/set"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	clusterAddr = envflag.String("CLUSTER", "", "The location of bootstrap address for clustering. If blank clustering will be disabled")
)

const (
	port = ":5051"
)

func main() {
	envflag.Parse()
	if *clusterAddr != "" {
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger := log.New(os.Stdout, "server: ", log.Ldate|log.Ltime)
	s := grpc.NewServer()
	server, err := server.NewServer(*clusterAddr, logger)
	if err != nil {
		logger.Fatalf("unable to create server: %v\n", err)
	}
	set.RegisterAddOnlyServer(s, server)
	go s.Serve(lis)

}
