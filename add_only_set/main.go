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
	vals := make(map[string]int)
	set := &AddOnlyImpl{vals}
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
	vals map[string]int
}

var _ s.AddOnlySet = &AddOnlyImpl{}

func (a *AddOnlyImpl) Add(s string, r *s.Result) error {
	a.vals[s]++
	fmt.Print("Added\n")
	for k := range a.vals {
		fmt.Printf("key %s %d\n", k, a.vals[k])
	}
	return nil
}
func (a *AddOnlyImpl) Show(s string, r *s.Result) error {
	var keys []interface{}
	for k := range a.vals {
		keys = append(keys, k)
	}
	r.Data = keys
	return nil
}
