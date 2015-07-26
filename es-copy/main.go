package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var (
	source, target, indice, mtype string
)

func init() {
	flag.StringVar(&source, "source", "http://172.16.13.230:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&target, "target", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "target indice")
	flag.StringVar(&indice, "indice", "fsmedia2", "or fsmedia")
	flag.StringVar(&mtype, "mtype", "media", "target type")
}

func main() {
	flag.Parse()
	src, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(source))
	panic_error(err)
	target, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(target))
	panic_error(err)

	panic_error(xiuxiu.EsCreateIfNotExist(target, indice))

	xiuxiu.EsMediaScan(src, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(target, em, indice)
	})
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia, indice string) {

	if _, err := client.Index().Index(indice).Type(mtype).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
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
