package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

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
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	directors := xiuxiu.EmCleanName(em.Director)
	actors := xiuxiu.EmCleanName(em.Actor)
	if xiuxiu.EsDebug {
		fmt.Println(directors, em.Director, actors, em.Actor)
		return
	}
	em.Directors = directors
	em.Actors = actors
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
