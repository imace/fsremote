package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

func to_islice(s []string) (v []interface{}) {
	for _, i := range s {
		v = append(v, i)
	}
	return
}
func terms_query(client *elastic.Client, s string) (*elastic.SearchResult, error) {
	q := elastic.NewTermsQuery("pkgName", to_islice(strings.Fields(s))...)

	results, err := client.Search().Index(xiuxiu.EsAppIndice).Types("app").Query(&q).From(0).Size(200).Do()
	return results, err
}

func package_select(s string) (v []xiuxiu.EsApp) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	result, err := terms_query(client, s)
	panic_error(err)

	for _, hit := range result.Hits.Hits {
		var em xiuxiu.EsApp
		if err := json.Unmarshal(*hit.Source, &em); err != nil {
			log.Println(err)
		} else {
			v = append(v, em)
		}
	}
	return
}
