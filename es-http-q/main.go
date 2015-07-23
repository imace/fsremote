package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"github.com/hearts.zhang/fsremote"
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
func es_query(q string) (v []fsremote.EsMedia) {
	result, err := script_score_query(client, q)
	panic_error(err)
	var ttyp fsremote.EsMedia
	for _, item := range result.Each(reflect.TypeOf(ttyp)) {
		if em, ok := item.(fsremote.EsMedia); ok {
			v = append(v, em)
		}
	}
	log.Println(result.TotalHits())
	return
}

/*
  "sort": [
    {
      "weight": {
        "order": "desc"
      }
    }
  ],*/
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
							"script": "_source.weight",
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
