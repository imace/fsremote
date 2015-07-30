package main

import "net/rpc"

func jieba_segment(text string) Terms {
	client, err := rpc.DialHTTP("tcp", jieba)
	panic_error(err)
	var terms Terms
	client.Call("Jieba.Segment", "小明硕士毕业于中国科学院计算所，后在日本京都大学深造", &terms.Terms)
	return terms
}
