package main

import (
	"fmt"
	"github.com/rayram23/crdt/addonly/set"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":5051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	set.RegisterAddOnlyServer(s, NewServer())
	s.Serve(lis)
}

type Server struct {
	data map[string]struct{}
}

var _ set.AddOnlyServer = &Server{}

func NewServer() *Server {
	return &Server{
		make(map[string]struct{}),
	}
}

func (s *Server) Add(c context.Context, r *set.AddRequest) (*set.BooleanResponse, error) {
	var t bool
	if _, ok := s.data[*r.Val]; ok {
		t = false
	} else {

		s.data[*r.Val] = struct{}{}
		t = true
	}
	fmt.Printf("called added: %s %t\n", *r.Val, t)
	return &set.BooleanResponse{Resp: &t}, nil
}
func (s *Server) Size(c context.Context, b *set.BlankMessage) (*set.IntResponse, error) {
	size := uint64(len(s.data))
	fmt.Printf("returning size:  %v\n", size)
	return &set.IntResponse{Resp: &size}, nil
}
func (s *Server) Contains(c context.Context, b *set.ContainsRequest) (*set.BooleanResponse, error) {
	_, ok := s.data[*b.Val]
	fmt.Printf("Contains :%s? %v\n", *b.Val, ok)
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
