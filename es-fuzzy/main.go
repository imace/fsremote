package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

var (
	addr   string
	_fuzzy = NewModel()
)

type handler func(w http.ResponseWriter, r *http.Request)

func init() {
	flag.StringVar(&addr, "addr", ":8089", "listen address")
}
func main() {
	flag.Parse()
	load_medias()

	rpcsvr := new(Fuzzy)
	rpc.Register(rpcsvr)
	rpc.HandleHTTP()

	log.Println("start server")
	http.Handle("/fuzzy/term", handler(handle_fuzzy_term)) //term=

	http.ListenAndServe(addr, nil)
}

func load_medias() {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)
	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(em)
	})
}
func add_fuzzy_words(words []string, weight int, hint string) {
	for _, word := range words {
		if len([]rune(word)) > 1 {
			_fuzzy.SetCount(word, weight, hint, true)
		}
	}
}

func when_es_media(m xiuxiu.EsMedia) {
	weight, id := int(m.Weight*100), strconv.Itoa(m.MediaID)

	add_fuzzy_words(m.NameNorm, weight, id)
	add_fuzzy_words(m.Actors, weight, id)
	add_fuzzy_words(m.Roles, weight, id)
	add_fuzzy_words(m.Directors, weight, id)
	add_fuzzy_words(strings.Fields(m.Tags), weight, id)
}

//term=
func handle_fuzzy_term(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	term := r.FormValue("term")
	x := fuzzy_suggest(term)
	panic_error(json.NewEncoder(w).Encode(map[string]interface{}{"items": x}))
}
func fuzzy_suggest(term string) (v []int) {
	flags := map[int]struct{}{}
	_, candis := _fuzzy.Suggestions(term, true)
	for _, c := range candis {
		if id, err := strconv.Atoi(c.Snippet); err == nil {
			if _, ok := flags[id]; !ok {
				v = append(v, id)
				flags[id] = struct{}{}
			}
		}
	}
	return
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func (imp handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()
	imp(w, r)
}
