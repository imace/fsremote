package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hearts.zhang/fsremote"
	"github.com/olivere/elastic"
)

var (
	es string
	q  string
)

const (
	indice = "fsmedia2"
	mtype  = "media"
)

func init() {
	flag.StringVar(&es, "es", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&q, "q", "first", "query tags")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es))
	panic_error(err)

	result, err := script_score_query(client, q)
	panic_error(err)

	log.Println(result.TotalHits())
}

/*
{
  "query": {
    "function_score": {
      "functions": [
        {
          "script_score": {
            "script": "_score * _source.weight"
          }
        }
      ],
      "query": {
        "term": {
          "tags": "first"
        }
      }
    }
  }
}*/

func script_score_query(client *elastic.Client, q string) (*elastic.SearchResult, error) {
	body := map[string]interface{}{
		"sort": []map[string]interface{}{
			map[string]interface{}{
				"weight": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"functions": []map[string]interface{}{
					map[string]interface{}{
						"script_score": map[string]interface{}{
							"script": "_score * doc['weight'].value",
						},
					},
				},
				"query": map[string]interface{}{
					"term": map[string]interface{}{
						"tags": q,
					},
				},
			},
		},
	}
	///fsmedia2/media/_search
	res, err := client.PerformRequest("POST", "/fsmedia2/media/_search", url.Values{}, body)
	if err != nil {
		return nil, err
	}

	// Return search results
	ret := new(elastic.SearchResult)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
func when_es_media(em fsremote.EsMedia) {
	em.Weight, em.Weight2 = fsremote.MediaScore(em.Day, em.Week, em.Seven, em.Month, em.Play, int64(em.Release), em.DisplayType)
	print_es_media(em)
}
func print_es_media(em fsremote.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), f2s(em.Weight2), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
