package main

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

func hans_pinyin(hans string) []string {
	duoyin := pinyin.Pinyin(hans, heteronym())

	pinyins1 := combine_pinyin(duoyin, "", true)
	pinyins2 := combine_pinyin(duoyin, "-", false)
	return uniq_string(append(pinyins1, pinyins2...))
}

func combine_pinyin(pinyins [][]string, sep string, soundex bool) (v []string) {
	if len(pinyins) == 0 {
		return
	}
	first := pinyins[0]
	left := pinyins[1:]

	leftc := combine_pinyin(left, sep, soundex)

	for _, item := range first {
		v = append(v, append_pinyins(item, leftc, sep, soundex)...)
	}
	return v
}

func append_pinyins(first string, left []string, sep string, soundex bool) (v []string) {
	firsts := pinyin_soundex(first, soundex)
	if len(left) == 0 {
		return firsts
	}
	for _, first := range firsts {
		for _, item := range left {
			if len(item) > 0 {
				item = sep + item
			}
			v = append(v, first+item)
		}
	}
	return v
}
func heteronym() pinyin.Args {
	a := pinyin.NewArgs()
	a.Heteronym = true
	return a
}

//声母模糊音：s <--> sh，c<-->ch，z <-->zh，l<-->n，f<-->h，r<-->l，
//韵母模糊音：an<-->ang，en<-->eng，in<-->ing，ian<-->iang，uan<-->uang。
func pinyin_soundex(pinyin string, soundex bool) (v []string) {
	v = append(v, pinyin)
	if !soundex {
		return
	}
	if p := pinyin_soundex_suffix(pinyin_soundex_prefix(pinyin)); p != pinyin {
		v = append(v, p)
	}

	return
}

func pinyin_soundex_prefix(pinyin string) (b string) {
	b = pinyin

	switch {
	case strings.HasPrefix(pinyin, "sh"):
		b = "s" + strings.TrimPrefix(pinyin, "sh")
	case strings.HasPrefix(pinyin, "ch"):
		b = "c" + strings.TrimPrefix(pinyin, "ch")
	case strings.HasPrefix(pinyin, "zh"):
		b = "z" + strings.TrimPrefix(pinyin, "zh")
	case strings.HasPrefix(pinyin, "n"):
		b = "l" + strings.TrimPrefix(pinyin, "n")
	case strings.HasPrefix(pinyin, "h"):
		b = "f" + strings.TrimPrefix(pinyin, "h")
	case strings.HasPrefix(pinyin, "l"):
		b = "n" + strings.TrimPrefix(pinyin, "l")
	}
	return
}
func pinyin_soundex_suffix(pinyin string) (b string) {
	b = pinyin
	switch {
	case strings.HasSuffix(pinyin, "uang"):
		b = strings.TrimSuffix(pinyin, "uang") + "uan"
	case strings.HasSuffix(pinyin, "iang"):
		b = strings.TrimSuffix(pinyin, "iang") + "ian"
	case strings.HasSuffix(pinyin, "ang"):
		b = strings.TrimSuffix(pinyin, "ang") + "an"
	case strings.HasSuffix(pinyin, "eng"):
		b = strings.TrimSuffix(pinyin, "eng") + "en"
	case strings.HasSuffix(pinyin, "ing"):
		b = strings.TrimSuffix(pinyin, "ing") + "in"
	}
	return
}

func uniq_string(strs []string) (v []string) {
	flag := map[string]struct{}{}
	for _, str := range strs {
		if _, ok := flag[str]; !ok {
			v = append(v, str)
			flag[str] = struct{}{}
		}
	}
	return
}
