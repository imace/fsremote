package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hearts.zhang/fsremote"
)

type FaceSuggests struct {
	Suggests []FaceSuggest `json:"suggests,omitempty"`
}
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
func face_trim(dup []FaceSuggest, accu bool) FaceSuggests {
	v := face_uniq(dup, accu)

	fill_media(v)
	return FaceSuggests{v}
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

func face_address() string {
	return "http://" + (*face) + "/face/suggest/"
}

func face_suggest_wrap(q string) []FaceSuggest {
	x := face_suggest(q)
	log.Println("face return", len(x))

	if len(x) < 5 {
		x = append(x, face_split_suggest(q)...)
		log.Println("fuzzy return", len(x))
	}
	if len(x) < 5 {
		x = append(x, es_search(q)...)
		log.Println("es return", len(x))
	}
	return x
}
func face_suggest(q string) []FaceSuggest {
	log.Println("search", q)
	params := url.Values{}
	params.Add("q", q)
	uri := face_address() + "?" + params.Encode()
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

func es_address() string {
	return "http://" + (*es) + "/search.php"
}

func es_search(q string) (v []FaceSuggest) {
	params := url.Values{}
	params.Add("tags", q)
	uri := es_address() + "?" + params.Encode()
	log.Println(uri)
	resp, err := http.Get(uri)
	panic_error(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic_error(errors.New("status not ok"))
	}

	xv := fsremote.EsMedias{}
	err = json.NewDecoder(resp.Body).Decode(&xv)
	for _, m := range xv.Data {
		item := fsremote.FunMedia{
			m.MediaID, m.Name, m.Name, m.NameEn, m.NameOt, m.Lang, m.MediaLength, m.Country, 0, m.CoverPicID, strings.Fields(m.Tags), 0,
		}

		sc, _ := strconv.Atoi(m.Release)
		v = append(v, FaceSuggest{"", score(m.Day, m.Week, m.Seven, m.Month, m.Play, sc), strconv.Itoa(m.MediaID), 0, item})
	}
	return
}

const _2020 = 1577836800

func score(d, w, s7, m, t, date int) int {
	r, l, dt := time.Unix(0, 0), time.Unix(_2020, 0), time.Unix(int64(date), 0)
	rg := math.Log(l.Sub(r).Hours() / 24)
	dr := dt.Sub(r).Hours() / 24
	dl := l.Sub(dt).Hours() / 24
	if dr < math.E {
		dr = math.E
	}
	if dl < math.E {
		dl = math.E
	}
	r1, l1 := math.Log(dr)/rg, math.Log(dl)/rg
	weight := r1*1.0 + l1*2.0
	days := time.Since(time.Unix(int64(date), 0)).Hours() / 24
	if days < 1.0 {
		days = 1.0
	}
	x := float64(t)/days + float64(d) + float64(s7)/7.0 + float64(w)/5.0 + float64(m)/30.0
	return int(x * weight)
}
