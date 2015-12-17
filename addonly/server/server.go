package server

import (
	"fmt"
	"github.com/hashicorp/serf/client"
	"github.com/rayram23/crdt/addonly/set"
	"golang.org/x/net/context"
	"log"
)

type Server struct {
	data        map[string]struct{}
	clusterAddr string
	rpcClient   *client.RPCClient
	logger      *log.Logger
}

var _ set.AddOnlyServer = &Server{}

func clusterUp(clusterAddr string) (*client.RPCClient, error) {
	client, err := client.NewRPCClient(clusterAddr)
	if err != nil {
		return nil, err
	}
	_, err = client.Join([]string{clusterAddr}, false)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewServer(clusterAddr string, logger *log.Logger) (*Server, error) {

	var client *client.RPCClient
	var err error
	if client, err = clusterUp(clusterAddr); err != nil {
		return nil, err
	}

	return &Server{
		data:        make(map[string]struct{}),
		clusterAddr: clusterAddr,
		rpcClient:   client,
		logger:      logger,
	}, nil
}

func (s *Server) Add(c context.Context, r *set.AddRequest) (*set.BooleanResponse, error) {
	var t bool
	if _, ok := s.data[*r.Val]; ok {
		t = false
	} else {
		s.data[*r.Val] = struct{}{}
		t = true
	}
	return &set.BooleanResponse{Resp: &t}, nil
}
func (s *Server) Size(c context.Context, b *set.BlankMessage) (*set.IntResponse, error) {
	size := uint64(len(s.data))
	return &set.IntResponse{Resp: &size}, nil
}
func (s *Server) Contains(c context.Context, b *set.ContainsRequest) (*set.BooleanResponse, error) {
	_, ok := s.data[*b.Val]
	return &set.BooleanResponse{Resp: &ok}, nil
}
func (s *Server) All(c context.Context, b *set.BlankMessage) (*set.AllResponse, error) {
	vals := make([]string, len(s.data))
	i := 0
	for k, _ := range s.data {
		vals[i] = k
		i++
	}
	fmt.Print("returned all\n")
	return &set.AllResponse{Val: vals}, nil
}
