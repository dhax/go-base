package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"

	"github.com/dhax/go-base/database"
	"github.com/dhax/go-base/logging"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
)

// API provides application resources and handlers.
type API struct {
	Account *AccountResource
}

// NewAPI configures and returns application API.
func NewAPI(db *pg.DB) (*API, error) {
	accountStore := database.NewAccountStore(db)
	account := NewAccountResource(accountStore)

	api := &API{
		Account: account,
	}
	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/account", a.Account.router())

	return r
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
