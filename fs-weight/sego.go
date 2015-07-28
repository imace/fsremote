package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
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

func init() {
	_stop_words = make(map[string]struct{})
	_stop_words["的"] = struct{}{}
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
