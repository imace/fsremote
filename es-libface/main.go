package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/hearts.zhang/fsremote"
	"github.com/olivere/elastic"
)

var (
	es      string
	tindice string
)

const (
	indice = "fsmedia2"
	mtype  = "media"
)

func init() {
	flag.StringVar(&es, "es", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&tindice, "indice", "fsmedia2", "target indice")
}

var (
	media fsremote.EsMedia
)

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es))
	panic_error(err)

	cursor, err := client.Scan(indice).Type(mtype).Size(100).Do()
	panic_error(err)

	for {
		result, err := cursor.Next()
		if err == elastic.EOS {
			break
		}
		panic_error(err)
		for _, hit := range result.Hits.Hits {
			var em fsremote.EsMedia
			if err := json.Unmarshal(*hit.Source, &em); err != nil {
				log.Println(err)
			} else {
				when_es_media(client, em)
			}
		}
	}

}

func uniq_string(a []string) (v []string) {
	set := make(map[string]struct{})
	for _, i := range a {
		set[i] = struct{}{}
	}
	for k := range set {
		v = append(v, k)
	}
	return v
}
func when_es_media(client *elastic.Client, em fsremote.EsMedia) {
	media = em
	var v = strings.Fields(em.Tags)
	for _, name := range em.NameNorm {
		v = append(v, norm_word(name)...)
	}
	v = append(v, norm_word(em.Director)...)
	v = append(v, norm_word(em.Actor)...)
	v = append(v, norm_word(em.Role)...)

	v = uniq_string(v)
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
