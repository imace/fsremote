package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"

	"github.com/huichen/sego"
)

var (
	addr, dict string
	_puncts    string = "的，在。？！、；：“” ‘’（）─…—·《》【】［］〈〉+-×÷≈＜＞%‰∞∝√∵∴∷∠⊙○π⊥∪°′〃℃{}()[].|&*/#~.,:;?!'-→．"
)

type Terms struct {
	Terms []string `json:"terms,omitempty"`
}
type handler func(w http.ResponseWriter, r *http.Request)

func init() {
	flag.StringVar(&addr, "addr", ":8081", "listen address")
	flag.StringVar(&dict, "dict", "e:/sego.dic,e:/dictionary.txt,e:/dict.txt", "sego dict, user dict first")
}
func main() {
	flag.Parse()

	_segmenter.LoadDictionary(dict)

	http.Handle("/sego", handler(handle_sego)) //?text=
	http.ListenAndServe(addr, nil)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

func is_stop_word(seg string) bool {
	return strings.Contains(_puncts, seg)
}
func handle_sego(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	segs := segment(text)
	var terms Terms
	for _, seg := range segs {
		if len(seg) > 1 && !is_stop_word(seg) {
			terms.Terms = append(terms.Terms, seg)
		}
	}
	panic_error(json.NewEncoder(w).Encode(&terms))
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

var _segmenter sego.Segmenter

func segment(text string) []string {
	v := _segmenter.Segment([]byte(text))
	return sego.SegmentsToSlice(v, true)
}
