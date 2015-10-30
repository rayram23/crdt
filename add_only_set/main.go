package main

import (
	s "./structs"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	fmt.Print("Hello World\n")
	set := &AddOnlyImpl{}
	rpc.Register(set)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		fmt.Printf("Error %s")
		return
	}
	http.Serve(l, nil)

}

type AddOnlyImpl struct {
	vals map[interface{}]int
}

var _ s.AddOnlySet = &AddOnlyImpl{}

func (a *AddOnlyImpl) Add(v interface{}, r *s.Result) error {
	a.vals[v]++
	fmt.Print("Added\n")
	return nil
}
func (a *AddOnlyImpl) Show(v interface{}, r *s.Result) error {
	var keys []interface{}
	for k := range a.vals {
		keys = append(keys, k)
	}
	r.Data = keys
	return nil
}
