package middlewares

import (
	"log"
	"net/http"
)

// Auth is an http middleware to authenticate http request
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("authenticating...")

		// TODO:
		// implement authentication middleware
		// ...

		next.ServeHTTP(w, r)
	})
}
