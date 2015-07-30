package main

import "github.com/huichen/sego"

type Sego sego.Segmenter

func (*Sego) Segment(txt string, terms *[]string) error {
	*terms = segment(txt)
	return nil
}
