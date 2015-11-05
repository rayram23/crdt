package addonly

import (
	"golang.org/x/net/context"
)

type Server struct {
	data map[string]struct{}
}

var _ AddOnlyServer = &Server{}

func NewServer() *Server {
	return &Server{
		make(map[string]struct{}),
	}
}

func (s *Server) Add(c context.Context, r *AddRequest) (*BooleanResponse, error) {
	s.data[*r.Val] = struct{}{}
	t := true
	return &BooleanResponse{Resp: &t}, nil
}
func (s *Server) Size(c context.Context, b *BlankMessage) (*IntResponse, error) {
	size := uint64(len(s.data))
	return &IntResponse{Resp: &size}, nil
}
func (s *Server) Contains(c context.Context, b *ContainsRequest) (*BooleanResponse, error) {
	_, ok := s.data[*b.Val]
	return &BooleanResponse{Resp: &ok}, nil
}
func (s *Server) All(c context.Context, b *BlankMessage) (*AllResponse, error) {
	vals := make([]string, len(s.data))
	i := 0
	for k, _ := range s.data {
		vals[i] = k
		i++
	}
	return &AllResponse{Val: vals}, nil
}
