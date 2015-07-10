package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/hearts.zhang/fsremote"
)

type FaceSuggest struct {
	Phrase  string            `json:"phrase"`
	Score   int               `json:"score"`
	Snippet string            `json:"snippet,omitempty"`
	Dup     int               `json:"dup"`
	Media   fsremote.FunMedia `json:"media,omitempty"`
}

func load_medias() {
	file, err := os.Open(*medias_file)
	panic_error(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		to_fun_media(line)
	}
}
func to_fun_media(line string) {
	var m fsremote.FunMedia
	panic_error(json.Unmarshal([]byte(line), &m))
	_medias[m.MediaId] = &m
	fill_rune2medias(&m)
}

func fill_rune2medias(m *fsremote.FunMedia) {
	_fuzzy.SetCount(m.Name, 1, strconv.Itoa(m.MediaId), true)
}

func face_uniq(dup []FaceSuggest, accu bool) (v []FaceSuggest) {
	x := make(map[string]*FaceSuggest)
	for _, suggest := range dup {
		if v, ok := x[suggest.Snippet]; ok {

			if accu {
				v.Score = v.Score + suggest.Score
				v.Dup++
			} else if v.Score < suggest.Score {
				var tmp = suggest
				x[suggest.Snippet] = &tmp
			}
		} else {
			var tmp = suggest
			x[suggest.Snippet] = &tmp
		}

	}
	var h = &FaceSuggestSlice{}
	heap.Init(h)
	for _, val := range x {
		heap.Push(h, *val)
	}
	for h.Len() > 0 && len(v) < limit {
		v = append(v, heap.Pop(h).(FaceSuggest))
	}
	return
}

type FaceSuggestSlice []FaceSuggest

func (slice FaceSuggestSlice) Len() int {
	return len(slice)
}

func (s FaceSuggestSlice) Less(i, j int) bool {
	if s[i].Dup > s[j].Dup {
		return true
	} else if s[i].Dup == s[j].Dup {
		return s[i].Score > s[j].Score
	} else {
		return false
	}
}

func (slice FaceSuggestSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (h *FaceSuggestSlice) Push(x interface{}) {
	*h = append(*h, x.(FaceSuggest))
}
func (h *FaceSuggestSlice) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//uniq and sort
func face_trim(dup []FaceSuggest, accu bool) []FaceSuggest {
	v := face_uniq(dup, accu)

	fill_media(v)
	return v
}

func fill_media(medias []FaceSuggest) {
	for i := 0; i < len(medias); i++ {
		medias[i].Media = *_medias[media_id(medias[i])]
	}
}

func media_id(fs FaceSuggest) int {
	v, _ := strconv.Atoi(fs.Snippet)
	return v
}

func face_suggest(q string) []FaceSuggest {
	log.Println("search", q)
	params := url.Values{}
	params.Add("q", q)
	uri := "http://" + (*face) + "/face/suggest/?" + params.Encode()
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
	_, datas := _fuzzy.Suggestions(q, true)
	for _, data := range datas {
		m := _medias[atoi(data.Snippet)]
		v = append(v, FaceSuggest{m.Name, int(m.Weight), strconv.Itoa(m.MediaId), 1, *m})
	}

	return v
}
