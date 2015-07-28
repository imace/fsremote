package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"reflect"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var (
	addr   string
	es     string
	client *elastic.Client
)

const (
	indice = "fsmedia2"
	mtype  = "media"
)

func init() {
	flag.StringVar(&addr, "addr", ":9205", "listening address")
	flag.StringVar(&es, "es", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "or http://testbox02.chinacloudapp.cn:9200")
}
func main() {
	flag.Parse()
	var err error
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(es))
	panic_error(err)

	http.Handle("/query", handler(handle_face)) //q

	http.ListenAndServe(addr, nil)

}
func es_query(q string) (v []xiuxiu.EsMedia) {
	result, err := script_score_query(client, q)
	panic_error(err)
	var ttyp xiuxiu.EsMedia
	for _, item := range result.Each(reflect.TypeOf(ttyp)) {
		if em, ok := item.(xiuxiu.EsMedia); ok {
			v = append(v, em)
		}
	}
	log.Println(result.TotalHits())
	return
}

func script_score_query(client *elastic.Client, s string) (*elastic.SearchResult, error) {
	//q := elastic.NewMoreLikeThisFieldQuery("tags", "Golang topic.")
	q := elastic.NewFunctionScoreQuery().
		Query(elastic.NewTermQuery("tags", s)).
		AddScoreFunc(elastic.NewFieldValueFactorFunction().Field("weight"))

	results, err := client.Search().Index("fsmedia2").Query(&q).From(0).Size(200).Do()
	return results, err
}
func handle_face(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.FormValue("q")
	x := es_query(q)
	panic_error(json.NewEncoder(w).Encode(x))
}

type handler func(w http.ResponseWriter, r *http.Request)

func (imp handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()
	imp(w, r)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
