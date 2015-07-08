// fs-import project main.go
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"sync"
)

type handler func(w http.ResponseWriter, r *http.Request)

var q = flag.String("q", "刘x德华", "query word")

func main() {
	flag.Parse()
	x := face_trim(face_suggest(*q), false)
	if len(x) == 0 {
		x = face_trim(face_split_suggest(*q), true)
	}
	for _, i := range x {
		fmt.Println(i)
	}
}
func main_test() {
	http.Handle("/face", handler(handle_face)) // ?lat=xxx&lng=xxx
	http.ListenAndServe(":9204", nil)
}

//[{ "phrase": "西亚特快", "score": 106, "snippet": "15627" }]
func handle_face(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.FormValue("q")
	x := face_trim(face_suggest(q), false)
	if len(x) == 0 {
		x = face_trim(face_split_suggest(q), true)
	}
	panic_error(json.NewEncoder(w).Encode(x))
}

type FaceSuggest struct {
	Phrase  string `json:"phrase"`
	Score   int    `json:"score"`
	Snippet string `json:"snippet,omitempty"`
}

func face_uniq(dup []FaceSuggest, accu bool) (v []FaceSuggest) {
	x := make(map[string]*FaceSuggest)
	for _, suggest := range dup {
		if v, ok := x[suggest.Snippet]; ok {
			if accu {
				v.Score = v.Score + suggest.Score
			} else if v.Score < suggest.Score {
				var tmp = suggest
				x[suggest.Snippet] = &tmp
			}
		} else {
			var tmp = suggest
			x[suggest.Snippet] = &tmp
		}
	}
	for _, val := range x {
		v = append(v, *val)
	}
	return
}

type FaceSuggestSlice []FaceSuggest

func (slice FaceSuggestSlice) Len() int {
	return len(slice)
}

func (slice FaceSuggestSlice) Less(i, j int) bool {
	return slice[i].Score > slice[j].Score
}

func (slice FaceSuggestSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//uniq and sort
func face_trim(dup []FaceSuggest, accu bool) []FaceSuggest {
	v := face_uniq(dup, accu)
	sort.Sort(FaceSuggestSlice(v))
	return v
}
func face_suggest(q string) []FaceSuggest {
	params := url.Values{}
	params.Add("q", q)
	uri := "http://testbox02.chinacloudapp.cn:9203/face/suggest/?" + params.Encode()
	resp, err := http.Get(uri)
	panic_error(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic_error(errors.New("status not ok"))
	}

	v := []FaceSuggest{}
	err = json.NewDecoder(resp.Body).Decode(&v)
	return v
}
func face_split_suggest(q string) []FaceSuggest {
	var v []FaceSuggest
	var wg sync.WaitGroup
	var rc = make(chan []FaceSuggest)
	for _, r := range []rune(q) {
		wg.Add(1)
		go func() {
			defer func() { recover() }() //?
			defer wg.Done()
			rc <- face_suggest(string(r))
		}()
	}
	go func() {
		wg.Wait()
		close(rc)
	}()
	for x := range rc {
		v = append(v, x...)
	}
	return v
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
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
