package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var (
	media xiuxiu.EsMedia
)

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})

}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	media = em
	var v = strings.Fields(em.Tags)

	append_phrase(v, int(em.Weight*10000), em.MediaID)
}

func print_words(v []string) {
	for _, word := range v {
		fmt.Println(word)
	}
}

func tag_suffix(tag string) (v []string) {
	v = append(v, tag)
	orig := []rune(tag)
	for i := 1; i < len(orig)-1; i++ {
		suffix := orig[i:]
		v = append(v, string(suffix))
	}
	return
}

func append_phrase(tags []string, weight, mediaid int) {
	for _, tag := range tags {
		append_imp(tag_suffix(tag), weight, mediaid)
	}
}

func append_imp(tags []string, weight, mediaid int) {
	for _, tag := range tags {
		fmt.Printf("%v\t%v\t%v\n", weight, tag, mediaid)
	}
}

func norm_word(word string) (v []string) {
	var current []rune
	var english bool
	for _, r := range []rune(word) {
		if !unicode.IsLetter(r) {
			if !english {
				v = push_noempty(v, &current)
			}
		} else if r < 256 {
			if !english {
				v = push_noempty(v, &current)
			} else {
				current = append(current, r)
			}
		} else {
			if english {
				current = nil
			} else {
				current = append(current, r)
			}
		}
	}
	if !english && len(current) > 0 {
		v = push_noempty(v, &current)
	}
	return v
}

func push_noempty(v []string, word *[]rune) []string {
	if len(*word) > 1 {
		v = append(v, string(*word))
	}

	*word = nil
	return v
}

func em_split_string(x string) (v []string) {
	fields := strings.Split(x, "/")
	for _, f := range fields {
		if x := strings.TrimSpace(f); x != "" {
			v = append(v, strings.TrimSpace(f))
		}
	}
	return
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

func f2s(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func drop_index(client *elastic.Client, index string) error {
	_, err := client.DeleteIndex(index).Do()
	return err
}

func create_index(client *elastic.Client, index string) error {
	_, err := client.CreateIndex(index).Do()
	return err
}
