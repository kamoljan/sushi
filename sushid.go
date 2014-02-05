package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	// "io"
	"io/ioutil"
	// "log"
	"net/http"
	"os"
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

func put(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		panic("Not supported Method")
	}

	b, err := ioutil.ReadAll(r.Body)

	h := sha1.New()
	h.Write(b)
	f := fmt.Sprintf("%x", h.Sum(nil))

	// write whole the body
	err = ioutil.WriteFile("/Users/kamol/work/go/gocode/src/github.com/kamoljan/sushi/"+f, b, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%s", b)
	check(err)
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "image-"+r.FormValue("id"))
}

func initStore(path string) {
	fmt.Println("Initializing data store...")
	for i := 0; i < 256; i++ {
		for x := 0; x < 256; x++ {
			err := os.MkdirAll(fmt.Sprintf("%s/%02x/%02x", path, i, x), 0755)
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("...Done") // total 65536 directories
}

func main() {
	// initialize data store
	initStore("/Users/kamol/work/go/gocode/src/github.com/kamoljan/sushi/store")

	http.HandleFunc("/", errorHandler(put))
	http.HandleFunc("/view", errorHandler(view))
	http.ListenAndServe(":8080", nil)

	// TEST
	// curl -XPUT http://localhost:8080/ad/saved -H "Content-type: image/jpeg" --data-binary @gopher.png
}
