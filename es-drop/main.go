package main

import (
	"flag"
	"fmt"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var (
	indice string
)

func init() {
	flag.StringVar(&indice, "dindice", "", "the dropping indice")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	if err != nil {
		panic(err)
	}
	if indice == "" {
		return
	}
	if xiuxiu.EsDebug {
		fmt.Println("drop", xiuxiu.EsAddr, indice)
	} else {
		drop_index(client, indice)
	}
}

func drop_index(client *elastic.Client, index string) error {
	_, err := client.DeleteIndex(index).Do()
	return err
}
