package main

import "github.com/huichen/sego"
import "github.com/hearts.zhang/xiuxiu/seg"

type Sego sego.Segmenter

func (*Sego) Segment(arg seg.SegoArg, terms *[]string) error {
	*terms = segment(arg.Text, !arg.IsNSearch)
	return nil
}
