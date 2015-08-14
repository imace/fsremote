package main

import (
	"flag"
	"fmt"
	"reflect"
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

	//apps := pkg_name_match("电视优化大师")
	apps := package_select("tv.fun.master", "电视优化大师")
	for _, app := range apps {
		when_es_media(app)
	}
}

func to_islice(s []string) (v []interface{}) {

	for _, i := range s {
		v = append(v, i)
	}
	return
}
func name_match_pkg_filter(client *elastic.Client, pkgs, name string) (*elastic.SearchResult, error) {
	f := elastic.NewTermsFilter("pkgName", to_islice(strings.Fields(pkgs))...)
	qq := elastic.NewMatchQuery("name", name)
	q := elastic.NewFilteredQuery(qq).Filter(f)
	results, err := client.Search().Index(xiuxiu.EsAppIndice).Types(xiuxiu.EsAppType).Query(&q).From(0).Size(200).Do()
	return results, err
}
func pkg_name_match(name string) (v []xiuxiu.EsApp) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)
	q := elastic.NewMatchQuery("name", name)
	results, err := client.Search().Index(xiuxiu.EsAppIndice).Types("app").Query(&q).From(0).Size(200).Do()
	panic_error(err)
	for _, iapp := range results.Each(reflect.TypeOf(xiuxiu.EsApp{})) {
		app := iapp.(xiuxiu.EsApp)
		v = append(v, app)
	}
	return
}
func package_select(pkgs, name string) (v []xiuxiu.EsApp) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	result, err := name_match_pkg_filter(client, pkgs, name)
	panic_error(err)
	for _, iapp := range result.Each(reflect.TypeOf(xiuxiu.EsApp{})) {
		app := iapp.(xiuxiu.EsApp)
		v = append(v, app)
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
	fmt.Println(em.PkgName)
}
