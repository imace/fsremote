// fs-import project main.go
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
)

const (
	edit_distance = 1
)

type Terms struct {
	Terms []string `json:"terms,omitempty"`
}
type handler func(w http.ResponseWriter, r *http.Request)

var (
	addr, sego, face, jieba string
	_medias                 = make(map[int]*xiuxiu.EsMedia)
	_fuzzy                  = NewModel()
)

func init() {
	flag.StringVar(&addr, "addr", ":8082", "listen address")
	flag.StringVar(&face, "face", "172.16.13.16:6767", "libface address")
	flag.StringVar(&sego, "sego", "172.16.13.16:8081", "sego address")
	flag.StringVar(&jieba, "jieba", "172.16.13.16:8083", "sego address")
	_fuzzy.SetDepth(edit_distance)
}

func main() {
	flag.Parse()
	load_medias()

	log.Println("start server")
	http.Handle("/app/select", handler(handle_app_select))
	http.Handle("/es/match", handler(handle_es_match))
	http.Handle("/fsmedia/face/term", handler(handle_face_term))   //t=term&n=
	http.Handle("/sego/seg", handler(handle_sego_seg))             //text=
	http.Handle("/jieba/seg", handler(handle_jieba_seg))           //text=
	http.Handle("/fsmedia/fuzzy/term", handler(handle_fuzzy_term)) //term=
	http.Handle("/img/sogou", handler(handle_img_sogou))           //q=&w=300&h=200
	http.Handle("/img/redirect.jpg", handler(handle_img_redirect)) //q=&w=200&h=400
	http.Handle("/pinyin/slug", handler(handle_pinyin_slug))       //hans=

	http.ListenAndServe(addr, nil)
}

//text=
func handle_jieba_seg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	terms := jieba_segment(text)
	panic_error(json.NewEncoder(w).Encode(&terms))
}

//text=
func handle_sego_seg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	text := r.FormValue("text")
	terms := sego_segment(text)
	panic_error(json.NewEncoder(w).Encode(&terms))
}

//term=
func handle_fuzzy_term(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	term := r.FormValue("term")
	x := fuzzy_trim(fuzzy_suggest(term))
	panic_error(json.NewEncoder(w).Encode(map[string]interface{}{"items": x}))
}

//t=term&n=
func handle_face_term(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	term := r.FormValue("t")
	n := atoi(r.FormValue("n"))
	if n < 1 {
		n = 16
	}
	x := face_trim(face_suggest(term, n))

	panic_error(json.NewEncoder(w).Encode(map[string]interface{}{"items": x}))
}

func handle_app_select(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
}

func handle_es_match(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

}

//q=&w=300&h=200
func handle_img_sogou(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q, pw, ph := r.FormValue("q"), r.FormValue("w"), r.FormValue("h")
	url, width, height := sogou_pic(q, atoi(pw), atoi(ph))

	panic_error(json.NewEncoder(w).Encode(map[string]interface{}{
		"uri":    url,
		"width":  width,
		"height": height,
	}))
}

//q=&w=300&h=200
func handle_img_redirect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.FormValue("q")
	q, pw, ph := r.FormValue("q"), r.FormValue("w"), r.FormValue("h")
	url, width, height := sogou_pic(q, atoi(pw), atoi(ph))
	w.Header().Del("Content-Type")
	w.Header().Set("Location", url)
	w.Header().Set("X-PIC", strconv.Itoa(width)+"x"+strconv.Itoa(height))
	w.WriteHeader(http.StatusFound)
}

//hans=
func handle_pinyin_slug(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hans := r.FormValue("hans")
	pinyin := hans_pinyin(hans)
	panic_error(json.NewEncoder(w).Encode(map[string]interface{}{
		"pinyin": pinyin,
	}))
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

func atoi(sid string) int {
	v, _ := strconv.Atoi(sid)
	return v
}
