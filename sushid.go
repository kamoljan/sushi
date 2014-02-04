package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Msg struct {
	Status, Message string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(500)
				w.Write(Message("ERROR", err.(string)))
			}
		}()
		fn(w, r)
	}
}

func Message(status string, message string) []byte {
	m := Msg{
		Status:  status,
		Message: message,
	}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err) // real panic
	}
	return b
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		panic("Not supported Method")
	}
	f, _, err := r.FormFile("image")
	check(err)
	defer f.Close()
	t, err := ioutil.TempFile("./", "image-")
	check(err)
	defer t.Close()
	_, err = io.Copy(t, f)
	check(err)
	http.Redirect(w, r, "/view?id="+t.Name()[6:], 302)
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "image-"+r.FormValue("id"))
}

func main() {
	http.HandleFunc("/", errorHandler(upload))
	http.HandleFunc("/view", errorHandler(view))
	http.ListenAndServe(":8080", nil)
}
