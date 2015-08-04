package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"reflect"
)

var (
	addr string
)

type handler func(w http.ResponseWriter, r *http.Request)

func init() {
	flag.StringVar(&addr, "addr", ":8082", "listen address")

}
func main() {
	flag.Parse()

	http.Handle("/semantic/hiv", handler(handle_semantic_hiv)) //name=&pkgs=
	http.ListenAndServe(addr, nil)
}

func from_map(s reflect.Value, v map[string]interface{}) {
	for k, m := range v {
		f := s.FieldByName(k)
		if f.IsValid() {
			tm := reflect.TypeOf(m)
			if tm.Kind() != reflect.Map {
				f.Set(reflect.ValueOf(m))
			} else if tm.Key().Kind() == reflect.String && tm.Elem().Kind() == reflect.Interface {
				from_map(f, m.(map[string]interface{}))
			}
		}
	}

}
func handle_semantic_hiv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := r.FormValue("s")
	sem := map[string]interface{}{}
	if s == "" {
		panic_error(json.NewDecoder(r.Body).Decode(&sem))
	} else {
		panic_error(json.Unmarshal([]byte(s), &sem))
	}
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
