package server

import (
	"log"
	"net/http"
)

func LoggedInOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Header.Get("Token"))
		log.Println(loggedAuthors[r.Header.Get("Token")])
		if _, ok := loggedAuthors[r.Header.Get("Token")]; !ok {
			http.NotFound(w, r)
			return
		}

		next(w, r)
	}
}
