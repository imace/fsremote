package main

import "testing"

func TestPinyin(t *testing.T) {
	x := hans_pinyin("长")
	t.Log(x)
}

func TestPinyinSoundex(t *testing.T) {
	t.Log(pinyin_soundex("chang", true))
	t.Log(pinyin_soundex("zhang", true))
}
