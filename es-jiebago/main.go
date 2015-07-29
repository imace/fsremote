package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/wangbin/jiebago"
)

var (
	addr                     string
	dict, userdict, segodict string
	_puncts                  string = "的，在。？！、；：“” ‘’（）─…—·《》【】［］〈〉+-×÷≈＜＞%‰∞∝√∵∴∷∠⊙○π⊥∪°′〃℃{}()[].|&*/#~.,:;?!'-→．"
	_segmenter               jiebago.Segmenter
)

type Terms struct {
	Terms []string `json:"terms,omitempty"`
}
type handler func(w http.ResponseWriter, r *http.Request)

func init() {
	flag.StringVar(&addr, "addr", ":8083", "listen address")
	flag.StringVar(&dict, "dict", "e:/dict.txt", "jieba default dict")
	flag.StringVar(&userdict, "user-dict", "e:/sego.dic", "media actors dict")
	flag.StringVar(&segodict, "sego-dict", "e:/dictionary.txt", "sego dictionary")
}
func main() {
	flag.Parse()

	_segmenter.LoadDictionary(dict)
	_segmenter.LoadUserDictionary(segodict)
	_segmenter.LoadUserDictionary(userdict)

	rpc.Register((*Jieba)(&_segmenter))
	rpc.HandleHTTP()

	http.Handle("/jieba", handler(handle_jieba)) //?text=
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
func segment(txt string) (v []string) {
	for term := range _segmenter.CutForSearch(txt, true) {
		if !is_stop_word(term) && len(term) > 1 {
			v = append(v, term)
		}
	}
	return
}
func handle_jieba(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	terms := Terms{segment(text)}

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
