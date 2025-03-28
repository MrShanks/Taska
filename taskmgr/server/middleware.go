package server

import (
	"log"
	"net/http"

	"github.com/MrShanks/Taska/utils"
)

func LoggedInOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if LoggedAuthorToken == "" {
			log.Println("Error: User not authenticated")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("You must login first")); err != nil {
				log.Printf("couldn't write response")
			}
			_, err := w.Write([]byte(nil))
			if err != nil {
				log.Printf("Couldn't write response: %v", err)
			}
			return
		}
		token, err := utils.VerifyToken(LoggedAuthorToken)
		if err != nil {
			log.Printf("Unauthorized: Token not valid or empty: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("You must login first")); err != nil {
				log.Printf("couldn't write response")
			}
			_, err := w.Write([]byte(nil))
			if err != nil {
				log.Printf("Couldn't write response: %v", err)
			}
			return
		}
		log.Printf("Login accomplished with token: %v\n", token)

		next(w, r)
	}
}
