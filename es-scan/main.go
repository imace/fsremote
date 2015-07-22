package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hearts.zhang/fsremote"
	"github.com/olivere/elastic"
)

var (
	es string
)

const (
	indice = "fsmedia"
	mtype  = "media"
)

func init() {
	flag.StringVar(&es, "es", "http://es.fun.tv:9200", "or http://testbox02.chinacloudapp.cn:9200")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es))
	panic_error(err)
	cursor, err := client.Scan(indice).Type(mtype).Size(20).Do()
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
			}

			when_es_media(em)
		}
	}
}
func when_es_media(em fsremote.EsMedia) {
	rl, _ := strconv.Atoi(em.Release)
	tl, base := time.Unix(int64(rl), 0), time.Unix(0, 0)
	delta := tl.Sub(base).Hours() / 24
	if delta < 1.0 {
		delta = 1.0
	}

	em.Weight, em.Weight2 = fsremote.MediaScore(em.Day, em.Week, em.Seven, em.Month, em.Play, int64(rl), em.DisplayType)
	print_es_media(em)
}
func print_es_media(em fsremote.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), f2s(em.Weight2), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
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
