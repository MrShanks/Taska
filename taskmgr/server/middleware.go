package server

import (
	"log"
	"net/http"

	"github.com/MrShanks/Taska/utils"
)

func LoggedInOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if utils.VerifyToken(LoggedAuthorToken) != nil {
			log.Println("Unauthorized: Token not valid or empty")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("You must login first")); err != nil {
				log.Printf("couldn't write response")
			}
			return
		}

		next(w, r)
	}
}
