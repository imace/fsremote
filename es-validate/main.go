package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

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
func append_misc(tags []string) (v []string, altered bool) {
	for _, item := range tags {
		v = append(v, item)
		switch item {
		case "文艺":
			altered = true
			v = append(v, "艺术")
		case "艺术":
			v = append(v, "文艺")
			altered = true
		}
	}
	return xiuxiu.EsUniqSlice(v), altered
}
func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	tags := strings.Fields(em.Tags)
	tags, x := append_misc(tags)
	if xiuxiu.EsDebug {
		if x {
			fmt.Println(em.Name, em.ReleaseDay)
		}
		return
	}
	if _, err := client.Index().Index(xiuxiu.EsIndice).Type(xiuxiu.EsType).Id(strconv.Itoa(em.MediaID)).BodyJson(&em).Do(); err != nil {
		log.Println(err)
	} else {
		fmt.Println(em.MediaID)
	}
}
func print_es_media(em xiuxiu.EsMedia) {
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
