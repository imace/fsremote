package main

import (
	"net/rpc"
	"testing"
)

func TestFuzyRpc(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", "localhost:8089")
	panic_error(err)
	var terms []int
	client.Call("Fuzzy.Guess", "刘德华", &terms)
	t.Log(terms)
}
