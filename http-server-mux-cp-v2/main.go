package main

import (
	"fmt"
	"net/http"
	"time"
)

func TimeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		weekday := now.Weekday().String()
		month := now.Month().String()
		day := now.Day()
		year := now.Year()
		result := fmt.Sprintf("%s, %d %s %d", weekday, day, month, year)
		w.Write([]byte(result))

	} // TODO: replace this
}

func SayHelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Hello there"
			w.Write([]byte(name))
		} else {
			result := fmt.Sprintf("Hello, %s!", name)
			w.Write([]byte(result))
		}

	} // TODO: replace this
}

func GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/time", TimeHandler())
	mux.HandleFunc("/hello", SayHelloHandler())
	// TODO: answer here
	return mux
}

func main() {
	http.ListenAndServe("localhost:8080", GetMux())
}
