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
	names := xiuxiu.EmCleanName(em.Name)
	names = append(names, xiuxiu.EmCleanName(em.NameEn)...)
	names = append(names, xiuxiu.EmCleanName(em.NameOt)...)
	names = xiuxiu.EsUniqSlice(names)
	if debug == true {
		fmt.Println(names, em.Name, em.NameEn, em.NameOt)
		return
	}
	em.NameNorm = names
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
