package main

import (
	"container/heap"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/olivere/elastic"
)

type FaceSuggests struct {
	Suggests []*xiuxiu.EsMedia `json:"suggests,omitempty"`
}
type FaceSuggest struct {
	Phrase  string `json:"phrase"`
	Score   int    `json:"score"`
	Snippet string `json:"snippet,omitempty"`
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
	_medias[m.MediaID] = &m
	weight, id := int(m.Weight*10000), strconv.Itoa(m.MediaID)
	//_fuzzy.SetCount(m.Name, int(m.Weight*10000), strconv.Itoa(m.MediaID), true)
	add_fuzzy_words(m.NameNorm, weight, id)
	add_fuzzy_words(m.Actors, weight, id)
	add_fuzzy_words(m.Roles, weight, id)
	add_fuzzy_words(m.Directors, weight, id)
	add_fuzzy_words(strings.Fields(m.Tags), weight, id)
}

func face_uniq(dup []FaceSuggest) (v []*xiuxiu.EsMedia) {
	x := map[int]struct{}{}
	for _, suggest := range dup {
		id := atoi(suggest.Snippet)
		if _, ok := x[id]; !ok {
			x[id] = struct{}{}
			v = append(v, _medias[id])
		}
	}
	return
}

type MediaHeap []*xiuxiu.EsMedia

func (slice MediaHeap) Len() int {
	return len(slice)
}

func (s MediaHeap) Less(i, j int) bool {
	return s[i].Weight > s[j].Weight
}

func (slice MediaHeap) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
func (h *MediaHeap) Push(x interface{}) {
	*h = append(*h, x.(*xiuxiu.EsMedia))
}
func (h *MediaHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type FaceSuggestSlice []FaceSuggest

func (slice FaceSuggestSlice) Len() int {
	return len(slice)
}

func (s FaceSuggestSlice) Less(i, j int) bool {
	return s[i].Score > s[j].Score
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
func face_trim(dup []FaceSuggest) FaceSuggests {
	v := face_uniq(dup)

	return FaceSuggests{v}
}

func media_id(fs FaceSuggest) int {
	v, _ := strconv.Atoi(fs.Snippet)
	return v
}

func face_address() string {
	return "http://" + face + "/face/suggest/"
}

func face_suggest(q string, n int) []FaceSuggest {
	log.Println("search", q)
	params := url.Values{}
	params.Add("q", q)
	params.Add("n", strconv.Itoa(n))
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
func fuzzy_suggest(term string) []TermData {
	_, v := _fuzzy.Suggestions(term, true)
	return v
}
func fuzzy_trim(v []TermData) (ret FaceSuggests) {
	medias := map[int]struct{}{}
	mh := &MediaHeap{}
	heap.Init(mh)
	for _, td := range v {
		id := atoi(td.Snippet)
		if _, ok := medias[id]; !ok {
			medias[id] = struct{}{}
			heap.Push(mh, _medias[id])
		}
	}
	count := 10
	if mh.Len() < 10 {
		count = mh.Len()
	} else if mh.Len() > 20 {
		count = 20
	}
	for i := 0; i < count; i++ {
		ret.Suggests = append(ret.Suggests, heap.Pop(mh).(*xiuxiu.EsMedia))
	}
	return
}
