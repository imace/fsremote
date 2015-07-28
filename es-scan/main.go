package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

func init() {

}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)
	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(em)
	})
}
func when_es_media(em xiuxiu.EsMedia) {
	print_es_media(em)
}
func print_es_media(em xiuxiu.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), em.MediaLength, unix_time(int(em.Release)), em.Tags)
}
func unix_time(t int) string {
	return time.Unix(int64(t), 0).Format(time.RFC3339)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
