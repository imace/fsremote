package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var (
	input string
)

func init() {
	flag.StringVar(&input, "input", "e:/words.txt", "query words")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
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
			var em xiuxiu.EsMedia
			if err := json.Unmarshal(*hit.Source, &em); err == nil {
				print_es_media(em)
			} else {
				log.Println(err)
			}
		}

	} else {
		log.Println(err)
	}
}

func print_es_media(em xiuxiu.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
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

func script_score_query(client *elastic.Client, s string) (*elastic.SearchResult, error) {
	//q := elastic.NewMoreLikeThisFieldQuery("tags", "Golang topic.")
	q := elastic.NewFunctionScoreQuery().
		Query(elastic.NewTermQuery("tags", s)).
		AddScoreFunc(elastic.NewFieldValueFactorFunction().Field("weight"))

	results, err := client.Search().Index("fsmedia2").Query(&q).From(0).Size(200).Do()
	return results, err
}
