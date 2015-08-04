package main

import (
	"net/url"
	"reflect"
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

func name_match_pkg_filter(client *elastic.Client, pkgs, name string) (*elastic.SearchResult, error) {
	//	f := elastic.NewTermsQuery("pkgName", to_islice(strings.Fields(s))...)
	f := elastic.NewTermsFilter("pkgName", to_islice(strings.Fields(pkgs))...)
	//q := elastic.NewTermQuery("name", "电视优化大师")
	//qq := elastic.NewQueryStringQuery("电视优化大师").DefaultField("name")
	qq := elastic.NewMatchQuery("name", name)
	q := elastic.NewFilteredQuery(qq).Filter(f)
	results, err := client.Search().Index(xiuxiu.EsAppIndice).Types(xiuxiu.EsAppType).Query(&q).From(0).Size(200).Do()
	return results, err
}

func package_name_match(name string) (v []xiuxiu.EsApp) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)
	//	f := elastic.NewTermsQuery("pkgName", to_islice(strings.Fields(s))...)
	//	f := elastic.NewTermsFilter("pkgName", to_islice(strings.Fields(pkgs))...)
	//q := elastic.NewTermQuery("name", "电视优化大师")
	//qq := elastic.NewQueryStringQuery("电视优化大师").DefaultField("name")
	q := elastic.NewMatchQuery("name", name)
	//	q := elastic.NewFilteredQuery(qq).Filter(f)
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

func es_app_select_url(name, pkgs string) string {
	//http://es.fun.tv/app/select?tags=%E6%B5%8F%E8%A7%88%E5%99%A8
	params := url.Values{}
	params.Add("tags", name)
	if len(pkgs) > 0 {
		params.Add("installed", pkgs)
	}
	return "http://" + es_front + "/app/select" + "?" + params.Encode()
}
