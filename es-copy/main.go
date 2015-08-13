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
	target, target_indice string
)

func init() {
	flag.StringVar(&target, "target", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "target")
	flag.StringVar(&target_indice, "tindice", "fsmedia4", "target indice")
}

func main() {
	flag.Parse()
	src, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)
	t, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(target))
	panic_error(err)

	if xiuxiu.EsDebug {
		fmt.Println("cp", xiuxiu.EsAddr, xiuxiu.EsIndice, target, target_indice)
		return
	}
	panic_error(xiuxiu.EsCreateIfNotExist(t, target_indice))

	xiuxiu.EsMediaScan(src, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(t, em, target_indice)
	})
}

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia, indice string) {

	if _, err := client.Index().Index(indice).Type(xiuxiu.EsType).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
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
