//
// jwtauth example
//
// Sample output:
//
// [peter@pak ~]$ curl -v http://localhost:3333/
// *   Trying ::1...
// * Connected to localhost (::1) port 3333 (#0)
// > GET / HTTP/1.1
// > Host: localhost:3333
// > User-Agent: curl/7.49.1
// > Accept: */*
// >
// < HTTP/1.1 200 OK
// < Date: Tue, 13 Sep 2016 15:53:17 GMT
// < Content-Length: 17
// < Content-Type: text/plain; charset=utf-8
// <
// * Connection #0 to host localhost left intact
// welcome anonymous%
//
//
// [peter@pak ~]$ curl -v http://localhost:3333/admin
// *   Trying ::1...
// * Connected to localhost (::1) port 3333 (#0)
// > GET /admin HTTP/1.1
// > Host: localhost:3333
// > User-Agent: curl/7.49.1
// > Accept: */*
// >
// < HTTP/1.1 401 Unauthorized
// < Content-Type: text/plain; charset=utf-8
// < X-Content-Type-Options: nosniff
// < Date: Tue, 13 Sep 2016 15:53:19 GMT
// < Content-Length: 13
// <
// Unauthorized
// * Connection #0 to host localhost left intact
//
//
// [peter@pak ~]$ curl -H"Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.PZLMJBT9OIVG2qgp9hQr685oVYFgRgWpcSPmNcw6y7M" -v http://localhost:3333/admin
// *   Trying ::1...
// * Connected to localhost (::1) port 3333 (#0)
// > GET /admin HTTP/1.1
// > Host: localhost:3333
// > User-Agent: curl/7.49.1
// > Accept: */*
// > Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.PZLMJBT9OIVG2qgp9hQr685oVYFgRgWpcSPmNcw6y7M
// >
// < HTTP/1.1 200 OK
// < Date: Tue, 13 Sep 2016 15:54:26 GMT
// < Content-Length: 22
// < Content-Type: text/plain; charset=utf-8
// <
// * Connection #0 to host localhost left intact
// protected area. hi 123%
//

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JwtAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwtauth.Claims{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {
	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, router())
}

func router() http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})
	})

	return r
}
