package app

import (
	"github.com/go-chi/chi"
	"github.com/go-pg/pg"

	"github.com/dhax/go-base/database"
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
