package main

import (
	"flag"
	"fmt"
	"github.com/rayram23/crdt/addonly/set"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var (
	addr   string
	action string
	val    string
)

func init() {
	flag.StringVar(&addr, "saddr", "localhost:5051", "The server location")
	flag.StringVar(&action, "action", "add", "action")
	flag.StringVar(&val, "value", "foobar", "value to add")
}
func main() {
	flag.Parse()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := set.NewAddOnlyClient(conn)
	if action == "add" && val != "" {
		add := &set.AddRequest{Val: &val}
		r, err := c.Add(context.Background(), add)
		if err != nil {
			fmt.Printf("Error %s", err.Error)
			return
		}
		fmt.Printf("Done adding: %t\n", r)
		return
	}
	if action == "contains" {
		contains := &set.ContainsRequest{Val: &val}
		r, err := c.Contains(context.Background(), contains)
		if err != nil {
			fmt.Printf("could not check contains: %s\n", err.Error())
			return
		}
		fmt.Printf("Contains: %s: %t\n", val, *r.Resp)
		return
	}
	if action == "size" {
		v := uint64(200)
		size := &set.BlankMessage{Time: &v}
		r, err := c.Size(context.Background(), size)
		if err != nil {
			fmt.Printf("could not check size %s\n", err.Error())
			return
		}
		fmt.Printf("Size is: %v\n", r.Resp)
		return
	}
	if action == "all" {
		v := uint64(200)
		all := &set.BlankMessage{Time: &v}
		r, err := c.All(context.Background(), all)
		if err != nil {
			fmt.Printf("Could not get all %s\n", err.Error())
			return
		}
		for i := 0; i < len(r.Val); i++ {
			fmt.Printf("val: %s\n", r.Val[i])
		}
		return
	}
	fmt.Printf("unknown action %s\n", action)
}
