package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var q string

func init() {
	flag.StringVar(&q, "q", "com.dianshiyouhua com.share.dcas", "query param")
}
func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	result, err := terms_query_pkg_filter(client, q, "电视优化大师")
	panic_error(err)

	for _, hit := range result.Hits.Hits {
		var em xiuxiu.EsApp
		if err := json.Unmarshal(*hit.Source, &em); err != nil {
			log.Println(err)
		} else {
			when_es_media(em)
		}
	}

	log.Println(result.TotalHits())
}

func to_islice(s []string) (v []interface{}) {

	for _, i := range s {
		v = append(v, i)
	}
	return
}
func terms_query_pkg_filter(client *elastic.Client, pkgs, name string) (*elastic.SearchResult, error) {
	//	f := elastic.NewTermsQuery("pkgName", to_islice(strings.Fields(s))...)
	f := elastic.NewTermsFilter("pkgName", to_islice(strings.Fields(pkgs))...)
	//q := elastic.NewTermQuery("name", "电视优化大师")
	//qq := elastic.NewQueryStringQuery("电视优化大师").DefaultField("name")
	qq := elastic.NewMatchQuery("name", name)
	q := elastic.NewFilteredQuery(qq).Filter(f)
	results, err := client.Search().Index(xiuxiu.EsAppIndice).Types("app").Query(&q).From(0).Size(200).Do()
	return results, err
}

func package_select(pkgs, name string) (v []xiuxiu.EsApp) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	result, err := terms_query_pkg_filter(client, pkgs, name)
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
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
func when_es_media(em xiuxiu.EsApp) {
	fmt.Println(em)
}
