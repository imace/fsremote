package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-engine.io"
)

func handle_eio(server *engineio.Server) {
	for {
		if conn, err := server.Accept(); err == nil {
			go handle_eio_conn(conn)
		} else {
			log.Println(err)
		}
	}
}

func handle_eio_conn(conn engineio.Conn) {
	log.Println("connected:", conn.Id())
	defer func() {
		conn.Close()
		log.Println("disconnected:", conn.Id())
	}()

	for {
		var err error
		t, msg, err := conn.NextReader()
		if err == nil {
			err = handle_eio_message(conn, t, msg)
		}
		if err != nil {
			return
		}
	}
}

func handle_eio_message(conn engineio.Conn, typ engineio.MessageType, msg io.ReadCloser) error {
	defer msg.Close()
	b, err := ioutil.ReadAll(msg)
	if err == nil {
		if typ == engineio.MessageText {
			err = handle_eio_string(conn, b)
		}
		if typ == engineio.MessageBinary {
			err = handle_eio_binary(conn, b)
		}
	}
	return err
}

func handle_eio_string(conn engineio.Conn, msg []byte) error {
	var v map[string]interface{}
	err := json.Unmarshal(msg, &v)
	log.Println(v)
	return err
}

func handle_eio_binary(conn engineio.Conn, msg []byte) error {
	return nil
}

var manager = newServerSessions()

func main() {
	server, err := engineio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.SetPingInterval(time.Second * 2)
	server.SetPingTimeout(time.Second * 3)
	server.SetSessionManager(manager)
	go handle_eio(server)

	http.Handle("/engine.io/", server)
	http.Handle("/semantic/push", handler(handle_semantic_push)) //q= POST&id=
	http.Handle("/connections/list", handler(handle_connections_list))
	log.Println("Serving at localhost:8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

type handler func(w http.ResponseWriter, r *http.Request)

func (imp handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()
	imp(w, r)
}

func handle_semantic_push(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var d []byte
	if q := r.FormValue("q"); q == "" {
		if xd, err := ioutil.ReadAll(r.Body); err == nil {
			d = xd
		}
	}
	id := r.FormValue("id")
	var ret int
	if len(d) > 0 {
		log.Println("broadcast sem ", len(d))
		ret = manager.foreach(func(conn engineio.Conn) {
			if id != "" && id != conn.Id() {
				return
			}
			if w, err := conn.NextWriter(engineio.MessageText); err == nil {
				defer w.Close()
				w.Write(d)
			}
		})
	}

	panic_error(json.NewEncoder(w).Encode(&ret))
}
func handle_connections_list(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var ids []string
	manager.foreach(func(conn engineio.Conn) {
		ids = append(ids, conn.Id())
	})
	panic_error(json.NewEncoder(w).Encode(map[string]interface{}{"ids": ids}))
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
