package main

//depends es-nameot/es-tags/es-typ-tags/es-digit
import (
	"flag"
	"fmt"
	"strings"
	"unicode"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})
	for term, info := range _terms {
		fmt.Println(term, info.weight, info.freq, info.weight/info.freq/info.freq+1, info.pos)
	}
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	print_words(em.NameNorm, int(em.Weight*1010), "nz")
	print_words(em.Actors, int(em.Weight*1000), "nr")
	print_words(em.Directors, int(em.Weight*800), "nr")
	print_words(em.Roles, int(em.Weight*700), "nr")
	tags := strings.Fields(em.Tags)
	print_words(tags, int(em.Weight*500), "n")
}

type Term struct {
	weight int
	freq   int
	pos    string
}

var _terms = map[string]*Term{}
var _puncts = "，。？！、；：“” ‘’（）─…—·《》〈〉+-×÷≈＜＞%‰∞∝√∵∴∷∠⊙○π⊥∪°′〃℃{}()[].|&*/#~.,:;?!'-→．"

func strip_space_punctuation(x string) string {
	var invalid bool

	for _, r := range []rune(x) {
		if unicode.IsDigit(r) || unicode.IsPunct(r) || strings.ContainsRune(_puncts, r) || (r < 256) {
			invalid = true
		}
	}
	if invalid {
		return ""
	}
	return x
}
func print_words(v []string, weight int, pos string) {
	for _, item := range v {
		item = strip_space_punctuation(item)
		if len(item) == 0 {
			continue
		}
		if term, ok := _terms[item]; ok {
			term.weight = term.weight + weight
			term.freq = term.freq + 1
		} else {
			_terms[item] = &Term{weight, 1, pos}
		}
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
