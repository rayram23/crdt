package main

import (
	s "./structs"
	"net/rpc"
	"testing"
)

func Test(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	reply := &s.Result{}
	err = client.Call("AddOnlyImpl.Add", "Ray", reply)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		return
	}

}
