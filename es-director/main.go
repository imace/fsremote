package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", true, "diagnose mode")
}

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.ESAddr))
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(client, em)
	})
}
func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	directors := xiuxiu.EmCleanDirector(em.Director)
	if debug == true {
		fmt.Println(directors, em.Director)
		return
	}
	em.Directors = directors
	if _, err := client.Index().Index(xiuxiu.EsIndice).Type(xiuxiu.EsType).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(em.MediaID)
	}
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
