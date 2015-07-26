// fs-import project main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hearts.zhang/xiuxiu"
)

const (
	limit         = 16
	edit_distance = 1
)

type handler func(w http.ResponseWriter, r *http.Request)

var (
	q           = flag.String("q", "剪刀x手x德华", "query word")
	addr        = flag.String("addr", ":9204", "listen address")
	face        = flag.String("face", "testbox02.chinacloudapp.cn:9203", "libface address")
	medias_file = flag.String("medias", "e:/medias.json", "media file")
	es          = flag.String("es", "es.fun.tv", "es search host")
)

var (
	_medias = make(map[int]*xiuxiu.FunMedia)
	_fuzzy  = NewModel()
)

const (
	hao_weather_api_key     = "b721bcdcf5ea4db78e1482fd2668a97c"
	hao_ip_location_api_key = "06ac888daa9e4c7b8add72845393c543"
	hao_ip_api_key          = "da3d89162b6f4bee94b11fa03e701522"
	hao_movie_api_key       = "27e7428ff5654317baf909d803927bb6"
	hao_video_api_key       = "9f374b640fb54a869c0a10a17d0a0103"
	hao_financial_api_key   = "d965c03cdd344b97b5df05786ef55279"
	hao_stock_api_key       = "f504cf57ca8248289ffa7aafd3e318b9"
	hao_tv_api_key          = "007b155e09b54447a8dc4c105c07057f"
)

func init() {
	_fuzzy.SetDepth(edit_distance)
}

func main_test() {
	flag.Parse()
	load_medias()

	x := face_trim(face_split_suggest(*q), true)

	for _, i := range x.Suggests {
		fmt.Println(i)
	}
}
func main() {
	flag.Parse()
	load_medias()
	log.Println("start server")
	http.Handle("/face", handler(handle_face)) //

	http.Handle("/app/select", handler(handle_app_select))
	http.Handle("/es/query", handler(handle_es_query)) //tags, q
	http.Handle("/es/match", handler(handle_es_match))
	http.ListenAndServe(*addr, nil)
}

//[{ "phrase": "西亚特快", "score": 106, "snippet": "15627" }]
func handle_face(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.FormValue("q")
	x := face_trim(face_suggest_wrap(q), false)

	panic_error(json.NewEncoder(w).Encode(x))
}

func handle_app_select(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
}

func handle_es_query(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	q := r.FormValue("q")
	x := face_trim(es_search(q), false)
	panic_error(json.NewEncoder(w).Encode(x))
}

func handle_es_match(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

}

//http://apis.haoservice.com/weather/ip?ip=202.108.250.241&key=b721bcdcf5ea4db78e1482fd2668a97c
func handle_weather_ip(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	panic_error(err)
	uri := hao_weather_ip(host)
	w.Header().Del("Content-Type")
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusFound)
}
func handle_weather_city(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	city := r.FormValue("city")
	log.Println(r.RemoteAddr, city)
	uri := hao_weather_city(city)
	w.Header().Del("Content-Type")
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusFound)
}
func handle_stock_hs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gid := r.FormValue("gid")
	uri := hao_stock_hs(gid)
	w.Header().Del("Content-Type")
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusFound)
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

func hao_weather_ip(ip string) string {
	return fmt.Sprintf("http://apis.haoservice.com/weather/ip?ip=%v&key=%v", ip, hao_weather_api_key)
}

//http://apis.haoservice.com/weather?cityname=北京&key=b721bcdcf5ea4db78e1482fd2668a97c
//http://apis.haoservice.com/lifeservice/stock/hs?gid=sh601009&key=f504cf57ca8248289ffa7aafd3e318b9
func hao_weather_city(cityname string) string {
	cn := url.QueryEscape(cityname)
	return fmt.Sprintf("http://apis.haoservice.com/weather?cityname=%v&key=%v", cn, hao_weather_api_key)
}
func hao_stock_hs(gid string) string {
	return fmt.Sprintf("http://apis.haoservice.com/lifeservice/stock/hs?gid=%v&key=%v", gid, hao_stock_api_key)
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
