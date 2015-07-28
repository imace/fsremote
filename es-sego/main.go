package main

import (
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"net/url"
	"strings"
)

var (
	addr, sego  string
	_stop_words map[string]struct{}
)

type Segments struct {
	Segments []struct {
		Text string `json:"text"`
		Pos  string `json:"pos"`
	} `json:"segments"`
}
type Terms struct {
	Terms []string `json:"terms,omitempty"`
}
type handler func(w http.ResponseWriter, r *http.Request)

func init() {
	flag.StringVar(&addr, "addr", ":8081", "listen address")
	flag.StringVar(&sego, "sego", "172.16.13.16:8080", "sego address")
	_stop_words = make(map[string]struct{})
	_stop_words["的"] = struct{}{}
}
func main() {
	flag.Parse()
	http.Handle("/sego", handler(handle_sego)) //?text=
	http.ListenAndServe(addr, nil)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

func is_stop_word(seg string) bool {
	_, ok := _stop_words[seg]
	return ok
}
func handle_sego(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	segs := segment(text)
	var terms Terms
	for _, seg := range segs.Segments {
		if strings.ContainsRune(seg.Pos, 'n') && !is_stop_word(seg.Text) {
			terms.Terms = append(terms.Terms, seg.Text)
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
func sego_address() string {
	return "http://" + (sego) + "/json"
}
func segment(text string) Segments {
	params := url.Values{}
	params.Add("text", text)

	uri := sego_address() + "?" + params.Encode()
	resp, err := http.Get(uri)
	panic_error(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic_error(errors.New("status not ok"))
	}

	var v Segments
	panic_error(json.NewDecoder(resp.Body).Decode(&v))
	return v
}
