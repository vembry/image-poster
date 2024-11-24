package middlewares

import (
	"log"
	"net/http"
)

// Auth is an http middleware to authenticate http request
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("authenticating...")

		// NOTE:
		// right now we bypass every single request's authentication
		// by defining "x-user-id" header directly on the request
		// i'll be using "x-user-id" to defined authenticated
		// user's id

		// TODO:
		// implement authentication middleware
		// ...

		next.ServeHTTP(w, r)
	})
}
