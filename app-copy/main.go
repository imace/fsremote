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
	source, target string
)

func init() {
	flag.StringVar(&source, "source", "http://172.16.13.230:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&target, "target", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "target elastic")

}

func main() {
	flag.Parse()
	src, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(source))
	panic_error(err)
	target, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(target))
	panic_error(err)

	panic_error(xiuxiu.EsCreateIfNotExist(target, xiuxiu.EsAppIndice))

	xiuxiu.EsAppScan(src, xiuxiu.EsAppIndice, xiuxiu.EsAppType, func(app xiuxiu.EsApp) {
		when_es_media(target, app, "ottpomme2")
	})
}

func when_es_media(client *elastic.Client, app xiuxiu.EsApp, indice string) {
	if _, err := client.Index().Index(indice).Type(xiuxiu.EsAppType).Id(strconv.Itoa(app.AppID)).BodyJson(&app).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(app.AppID)
	}
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
