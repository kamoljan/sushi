package main

import (
	"net/http"
	"regexp"
)

func getHandler(w http.ResponseWriter, r *http.Request, fid string) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write([]byte("This is GET request"))
}

var validPath = regexp.MustCompile("^/(fid)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/", makeHandler(getHandler))
	http.ListenAndServe(":8080", nil)
}
