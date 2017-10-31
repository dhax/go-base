// Package admin ties together administration resources and handlers.
package admin

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg"

	"github.com/dhax/go-base/auth/authorize"
	"github.com/dhax/go-base/database"
	"github.com/dhax/go-base/logging"
)

const (
	roleAdmin = "admin"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
)

// API provides admin application resources and handlers.
type API struct {
	Accounts *AccountResource
}

// NewAPI configures and returns admin application API.
func NewAPI(db *pg.DB) (*API, error) {

	accountStore := database.NewAdmAccountStore(db)
	accounts := NewAccountResource(accountStore)

	api := &API{
		Accounts: accounts,
	}
	return api, nil
}

// Router provides admin application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(authorize.RequiresRole(roleAdmin))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Admin"))
	})

	r.Mount("/accounts", a.Accounts.router())
	return r
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
