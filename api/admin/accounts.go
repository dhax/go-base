package admin

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-ozzo/ozzo-validation"

	"github.com/dhax/go-base/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// The list of error types returned from account resource.
var (
	ErrAccountValidation = errors.New("account validation error")
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
)

// AccountStore defines database operations for account management.
type AccountStore interface {
	List(f models.AccountFilter) (*[]models.Account, int, error)
	Create(*models.Account) error
	Get(id int) (*models.Account, error)
	Update(*models.Account) error
	Delete(*models.Account) error
}

// AccountResource implements account managment handler.
type AccountResource struct {
	Store AccountStore
}

// NewAccountResource creates and returns an account resource.
func NewAccountResource(store AccountStore) *AccountResource {
	return &AccountResource{
		Store: store,
	}
}

func (rs *AccountResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Post("/", rs.create)
	r.Route("/{accountID}", func(r chi.Router) {
		r.Use(rs.accountCtx)
		r.Get("/", rs.get)
		r.Put("/", rs.update)
		r.Delete("/", rs.delete)
	})
	return r
}

func (rs *AccountResource) accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "accountID"))
		if err != nil {
			render.Render(w, r, ErrBadRequest)
			return
		}
		account, err := rs.Store.Get(id)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxAccount, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type accountRequest struct {
	*models.Account
}

func (d *accountRequest) Bind(r *http.Request) error {
	return nil
}

type accountResponse struct {
	*models.Account
}

func newAccountResponse(a *models.Account) *accountResponse {
	resp := &accountResponse{Account: a}
	return resp
}

type accountListResponse struct {
	Accounts *[]models.Account `json:"accounts"`
	Count    int               `json:"count"`
}

func newAccountListResponse(a *[]models.Account, count int) *accountListResponse {
	resp := &accountListResponse{
		Accounts: a,
		Count:    count,
	}
	return resp
}

func (rs *AccountResource) list(w http.ResponseWriter, r *http.Request) {
	f := models.NewAccountFilter(r.URL.Query())
	al, count, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newAccountListResponse(al, count))
}

func (rs *AccountResource) create(w http.ResponseWriter, r *http.Request) {
	data := &accountRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Respond(w, r, ErrInvalidRequest(err))
		return
	}

	acc := data.Account
	if err := rs.Store.Create(acc); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrAccountValidation, err))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) get(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)
	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) update(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)
	data := &accountRequest{Account: acc}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	acc = data.Account
	if err := rs.Store.Update(acc); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrAccountValidation, err))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) delete(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*models.Account)
	if err := rs.Store.Delete(acc); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}
