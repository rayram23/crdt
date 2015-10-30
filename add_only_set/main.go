package main

import (
	"./structs"
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
	go http.Serve(l, nil)

}

type AddOnlyImpl struct {
	vals map[interface{}]interface{}
}

var _ AddOnlySet = &AddOnlyImpl{}

func (a *AddOnlyImpl) Add(v interface{}, r *string) error {
	fmt.Print("Got to add\n")
	return nil
}
func (a *AddOnlyImpl) Show(v interface{}, r []interface{}) error {
	fmt.Print("Got to show\n")
	return nil
}
