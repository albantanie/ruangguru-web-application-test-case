package main

import (
	"net/http"
)

func AdminHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Admin page"))
	}
}

func RequestMethodGetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method is not allowed"))
		}
	}) // TODO: replace this
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("role")
		if role == "ADMIN" {

			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(401)
			w.Write([]byte("Role not authorized"))
		}
	}) // TODO: replace this
}

func main() {
	// TODO: answer here
	http.Handle("/admin", RequestMethodGetMiddleware(AdminMiddleware(http.HandlerFunc(AdminHandler()))))
	http.ListenAndServe("localhost:8080", nil)
}
