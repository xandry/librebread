package tinkoff

import (
	"log"
	"net/http"
)

func CheckApplicationJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			log.Printf("Tinkoff: unresolved header")
			return
		}

		next.ServeHTTP(w, r)
	})
}
