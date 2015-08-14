package main

import "github.com/huichen/sego"
import "github.com/hearts.zhang/xiuxiu/seg"

type Sego sego.Segmenter

func (*Sego) Segment(arg seg.SegoArg, terms *[]string) error {
	*terms = segment(arg.Text, !arg.IsNSearch)
	return nil
}
func segment(text string, search_mode bool) (v []string) {
	vs := _segmenter.Segment([]byte(text))

	for _, seg := range sego.SegmentsToSlice(vs, search_mode) {
		if len(seg) > 1 && !is_stop_word(seg) {
			v = append(v, seg)
		}
	}
	return v
}

func segments(text string, search_mode bool) (v []seg.SegoToken) {
	segs := _segmenter.Segment([]byte(text))
	if search_mode {
		for _, s := range segs {
			v = append(v, extract_seg_tokens(s.Token())...)
		}
	} else {
		for _, s := range segs {
			v = append(v, to_seg_token(s.Token()))
		}
	}
	return
}

func extract_seg_tokens(t *sego.Token) (v []seg.SegoToken) {
	v = append(v, to_seg_token(t))
	for _, s := range t.Segments() {
		v = append(v, extract_seg_tokens(s.Token())...)
	}
	return v
}

func to_seg_token(t *sego.Token) (v seg.SegoToken) {
	v.Text = t.Text()
	v.Frequency = t.Frequency()
	v.Pos = t.Pos()
	return
}
