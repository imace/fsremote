package main

//depends es-nameot/es-tags/es-typ-tags/es-digit
import (
	"flag"
	"fmt"
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

func when_es_media(client *elastic.Client, em xiuxiu.EsMedia) {
	tags := strings.Fields(em.Tags)
	print_words(tags)
}
func print_words(v []string) {
	for _, word := range v {
		if len(word) > 0 && word[0] > 128 {
			fmt.Println(word)
		}
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
