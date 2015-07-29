package main

import "github.com/wangbin/jiebago"

type Jieba jiebago.Segmenter

func (*Jieba) Segment(txt string, terms *[]string) error {
	*terms = segment(txt)
	return nil
}
