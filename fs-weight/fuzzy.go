package main

import "net/rpc"

func fuzzy_suggest(text string) (v []int) {
	client, err := rpc.DialHTTP("tcp", fuzzy)
	panic_error(err)

	client.Call("Sego.Segment", "小明硕士毕业于中国科学院计算所，后在日本京都大学深造", &v)

	return
}
