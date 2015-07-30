package main

import (
	"net/rpc"
	"testing"
)

func TestSegoRpc(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	panic_error(err)
	var terms []string
	client.Call("Sego.Segment", "小明硕士毕业于中国科学院计算所，后在日本京都大学深造", &terms)
	for _, term := range terms {
		t.Log(term)
	}
}
