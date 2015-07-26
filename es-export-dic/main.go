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
	media fsremote.EsMedia
)

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	cursor, err := client.Scan(xiuxiu.EsIndice).Type(xiuxiu.EsType).Size(100).Do()
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
	print_special()
}
func print_special() {
	fmt.Println("电视")
	fmt.Println("电视剧")
	fmt.Println("系列剧")
	fmt.Println("连续剧")
	fmt.Println("电影")
	fmt.Println("电影片")
	fmt.Println("片儿")
	fmt.Println("卡通")
	fmt.Println("动漫")
	fmt.Println("漫画")
	fmt.Println("纪录片")
	fmt.Println("微电影")
	fmt.Println("微视频")
	fmt.Println("视频")
	fmt.Println("小视频")
	fmt.Println("微电视")
	fmt.Println("美剧")
	fmt.Println("英剧")
	fmt.Println("好莱坞")
	fmt.Println("日剧")
	fmt.Println("日韩")
	fmt.Println("韩剧")
	fmt.Println("泰剧")
	fmt.Println("港剧")
	fmt.Println("港片儿")
	fmt.Println("综艺")
	fmt.Println("娱乐")

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
	var v []string
	for _, name := range em.NameNorm {
		v = append(v, norm_word(name)...)
	}
	v = append(v, norm_word(em.Director)...)
	v = append(v, norm_word(em.Actor)...)
	v = append(v, norm_word(em.Role)...)
	v = uniq_string(v)
	print_words(v)
}
func print_words(v []string) {
	for _, word := range v {
		fmt.Println(word)
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
	if len(*word) <= 1 && len(*word) > 0 {
		log.Println(media.MediaID, media.Name, string(*word))
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
func print_es_media(em fsremote.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
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
