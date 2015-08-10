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

	http.Handle("/semantic/hiv", handler(handle_protocol_hiv))         //query= application/json
	http.Handle("/protocol/hiv", handler(handle_session_protocol_hiv)) //query= application/json
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
func handle_session_protocol_hiv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	q, s := r.FormValue("q"), r.FormValue("s")
	if q == "" {
		q = r.FormValue("query")
	}
	var v hi_session_protocol
	if s == "" && (r.Method == "POST" || r.Method == "PUT") {
		panic_error(json.NewDecoder(r.Body).Decode(&v))
	} else {
		panic_error(json.Unmarshal([]byte(s), &v))
	}
	v.text = q
	ret := when_session_protocol(v)
	panic_error(json.NewEncoder(w).Encode(ret))
}

func handle_protocol_hiv(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	q, s := r.FormValue("q"), r.FormValue("s")
	if q == "" {
		q = r.FormValue("query")
	}
	var v hi_protocol
	if s == "" && (r.Method == "POST" || r.Method == "PUT") {
		panic_error(json.NewDecoder(r.Body).Decode(&v))
	} else {
		panic_error(json.Unmarshal([]byte(s), &v))
	}

	if q != "" {
		v.Text = q
	}

	ret := when_protocol(v)
	panic_error(json.NewEncoder(w).Encode(ret))
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
