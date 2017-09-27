package api

import (
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/dhax/go-base/api/admin"
	"github.com/dhax/go-base/api/app"
	"github.com/dhax/go-base/auth"
	"github.com/dhax/go-base/database"
	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// NewAPI configures application resources and routes
func NewAPI() (*chi.Mux, error) {
	logger := logging.NewLogger()

	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}

	mailer, err := email.NewMailer()
	if err != nil {
		return nil, err
	}

	authStore := database.NewAuthStore(db)
	authResource, err := auth.NewResource(authStore, mailer)
	if err != nil {
		return nil, err
	}

	adminAPI, err := admin.NewAPI(db)
	if err != nil {
		return nil, err
	}

	appAPI, err := app.NewAPI(db)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// use CORS middleware if client is not served by this api, e.g. from other domain or CDN
	// r.Use(corsConfig().Handler)

	r.Mount("/auth", authResource.Router())
	r.Group(func(r chi.Router) {
		r.Use(authResource.Token.Verifier())
		r.Use(auth.Authenticator)
		r.Mount("/admin", adminAPI.Router())
		r.Mount("/api", appAPI.Router())
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	client := "./public"
	r.Get("/*", SPAHandler(client))

	return r, nil
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}

// SPAHandler serves the public Single Page Application
func SPAHandler(publicDir string) http.HandlerFunc {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()

		// serve static files
		if strings.Contains(url, ".") || url == "/" {
			handler.ServeHTTP(w, r)
			return
		}

		// otherwise always serve index.html
		http.ServeFile(w, r, path.Join(publicDir, "/index.html"))
	})
}
