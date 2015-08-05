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
	flag.StringVar(&addr, "addr", ":8086", "listen address")

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
func v_from_map(v interface{}, m map[string]interface{}) {
	val := reflect.ValueOf(v).Elem()
	from_map(val, m)
}
func handle_semantic_hiv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	q, s := r.FormValue("q"), r.FormValue("s")
	var v hi_understand_result
	if s == "" && (r.Method == "POST" || r.Method == "PUT") {
		panic_error(json.NewDecoder(r.Body).Decode(&v))
	} else {
		panic_error(json.Unmarshal([]byte(s), &v))
	}

	if q != "" {
		v.Text = q
	}

	when_understand(v)
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
