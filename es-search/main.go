package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var q string

func init() {
	flag.StringVar(&q, "q", "刘德华", "query param")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	result, err := script_score_query(client, q)
	panic_error(err)

	for _, hit := range result.Hits.Hits {
		var em xiuxiu.EsMedia
		if err := json.Unmarshal(*hit.Source, &em); err != nil {
			log.Println(err)
		} else {
			when_es_media(em)
		}
	}

	log.Println(result.TotalHits())
}

func script_score_query(client *elastic.Client, s string) (*elastic.SearchResult, error) {
	//q := elastic.NewMoreLikeThisFieldQuery("tags", "Golang topic.")
	q := elastic.NewFunctionScoreQuery().
		Query(elastic.NewTermQuery("tags", s)).
		AddScoreFunc(elastic.NewFieldValueFactorFunction().Field("weight"))

	results, err := client.Search().Index("fsmedia2").Query(&q).From(0).Size(200).Do()
	return results, err
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func when_es_media(em xiuxiu.EsMedia) {

	fmt.Println(em.MediaID)

}
