package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/neurodrone/crdt-1"
	"google.golang.org/grpc"
)

func main() {
	var (
		port     = flag.Int("port", -1, "port to start crdt server on")
		serverId = flag.String("server-id", "", "unique id for the server")
		peers    = flag.String("peers", "", "host:port comma-separated pairs of peers")
		timeout  = flag.Duration("timeout", 30*time.Second, "timeout for query")
	)
	flag.Parse()

	if *serverId == "" {
		log.Fatalln("a unique server id needs to be defined")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("cannot bind to port '%d': %s", *port, err)
	}

	grpcSrv := grpc.NewServer()
	srv := crdt.NewReplicationServer(*serverId, *timeout)

	crdt.RegisterReplicationTransportServer(grpcSrv, srv)

	go func() {
		peersList := make([]string, 1)
		peersList[0] = *peers
		pool, err := crdt.NewPeerPool(peersList, grpc.WithInsecure())
		if err != nil {
			log.Printf("cannot create peerpool:", err)
			return
		}

		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			if err := pingServers(pool); err != nil {
				log.Println("failed querying:", err)
			}
		}
	}()

	if err := grpcSrv.Serve(ln); err != nil {
		log.Fatalln("cannot start grpc server:", err)
	}
}

func pingServers(pool *crdt.PeerPool) error {
	for _, client := range pool.Clients() {
		if err := client.Ping(); err != nil {
			log.Println("ping failed:", err)
		} else {
			log.Println("ping succeeded to", client.Id())
		}
	}

	return nil
}
