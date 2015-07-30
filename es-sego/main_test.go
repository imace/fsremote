package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/huichen/sego"
)

func print(terms []string) {
	for _, word := range terms {
		fmt.Printf(" %s ", word)
	}
	fmt.Println()
}
func TestSego(t *testing.T) {
  t.Skip()
	segs := _segmenter.Segment([]byte("小明硕士毕业于中国科学院计算所，后在日本京都大学深造"))
	print(sego.SegmentsToSlice(segs, true))
}
func _TestMain(m *testing.M) {
	load_dicts()
	r := m.Run()
	os.Exit(r)
}
func load_dicts() {
	_segmenter.LoadDictionary("e:/sego.dic,e:/dictionary.txt,e:/dict.txt")
}
