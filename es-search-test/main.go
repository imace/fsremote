package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hearts.zhang/fsremote"
	"github.com/olivere/elastic"
)

var (
	es    string
	input string
)

const (
	indice = "fsmedia2"
	mtype  = "media"
)

func init() {
	flag.StringVar(&es, "es", "http://172.16.13.16:9200/", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&input, "input", "e:/words.txt", "query words")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es))
	panic_error(err)

	for_each_line(input, do_query, client)

}

func for_each_line(input string, do_query func(string, *elastic.Client), client *elastic.Client) {
	file, err := os.Open(input)
	panic_error(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		do_query(strings.TrimSpace(line), client)
	}
}
func do_query(q string, client *elastic.Client) {
	fmt.Println("search", q)
	if result, err := script_score_query(client, q); err == nil {
		for _, hit := range result.Hits.Hits {
			var em fsremote.EsMedia
			if err := json.Unmarshal(*hit.Source, &em); err == nil {
				fmt.Println(em.MediaID, em.Name, em.Weight)
			} else {
				log.Println(err)
			}
		}

	} else {
		log.Println(err)
	}
}
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
