package main

import "testing"

func TestSogouPic(t *testing.T) {
	t.Skip()
	t.Log(sogou_pic("张学友", 2, 1))
	t.Log(sogou_pic("小杨幂", 2, 1))
}
