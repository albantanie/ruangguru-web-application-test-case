package main

import (
	"fmt"
	"net/http"
	"time"
)

func TimeHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		weekday := now.Weekday().String()
		day := now.Day()
		month := now.Month().String()
		year := now.Year()
		result := fmt.Sprintf("%s, %d %s %d", weekday, day, month, year)
		w.Write([]byte(result))
	}) // TODO: replace this
}

func SayHelloHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name != "" {
			w.Write([]byte("Hello, " + name + "!"))
			return
		}
		w.Write([]byte("Hello there"))
	}) // TODO: replace this
}

func main() {
	// TODO: answer here

	http.HandleFunc("/time", TimeHandler())
	http.HandleFunc("/sayhello", SayHelloHandler())

	http.ListenAndServe("localhost:8080", nil)
}
