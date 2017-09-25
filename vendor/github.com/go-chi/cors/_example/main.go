// cors example
//
// ie.
//
// Unsuccessful Preflight request:
// ===============================
// $ curl -i http://localhost:3000/ -H "Origin: http://no.com" -H "Access-Control-Request-Method: GET" -X OPTIONS
// HTTP/1.1 200 OK
// Vary: Origin
// Vary: Access-Control-Request-Method
// Vary: Access-Control-Request-Headers
// Date: Fri, 28 Jul 2017 17:55:47 GMT
// Content-Length: 0
// Content-Type: text/plain; charset=utf-8
//
//
// Successful Preflight request:
// =============================
// $ curl -i http://localhost:3000/ -H "Origin: http://example.com" -H "Access-Control-Request-Method: GET" -X OPTIONS
// HTTP/1.1 200 OK
// Access-Control-Allow-Credentials: true
// Access-Control-Allow-Methods: GET
// Access-Control-Allow-Origin: http://example.com
// Access-Control-Max-Age: 300
// Vary: Origin
// Vary: Access-Control-Request-Method
// Vary: Access-Control-Request-Headers
// Date: Fri, 28 Jul 2017 17:56:44 GMT
// Content-Length: 0
// Content-Type: text/plain; charset=utf-8
//
//
// Content request (after a successful preflight):
// ===============================================
// $ curl -i http://localhost:3000/ -H "Origin: http://example.com"
// HTTP/1.1 200 OK
// Access-Control-Allow-Credentials: true
// Access-Control-Allow-Origin: http://example.com
// Access-Control-Expose-Headers: Link
// Vary: Origin
// Date: Fri, 28 Jul 2017 17:57:52 GMT
// Content-Length: 7
// Content-Type: text/plain; charset=utf-8
//
// welcome%
//
package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	http.ListenAndServe(":3000", r)
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	if origin == "http://example.com" {
		return true
	}

	return false
}
