package main

import (
	"fmt"
	"os"
	"testing"
)

func print(ch <-chan string) {
	for word := range ch {
		fmt.Printf(" %s /", word)
	}
	fmt.Println()
}
func TestJieba(t *testing.T) {
	t.Skip()
	print(_segmenter.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))
}
func _TestMain(m *testing.M) {

	load_dicts()
	r := m.Run()
	os.Exit(r)
}
func load_dicts() {
	_segmenter.LoadDictionary("e:/dict.txt")
	_segmenter.LoadUserDictionary("e:/dictionary.txt")
	_segmenter.LoadUserDictionary("e:/sego.dic")
}
