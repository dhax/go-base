package app

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/dhax/go-base/auth"
	"github.com/dhax/go-base/logging"
	"github.com/dhax/go-base/models"
)

// AccountStore defines database operations for account.
type AccountStore interface {
	Get(id int) (*auth.Account, error)
	Update(*auth.Account) error
	Delete(*auth.Account) error
	UpdateToken(*auth.Token) error
	DeleteToken(*auth.Token) error
	UpdateProfile(*models.Profile) error
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
	r.Use(rs.accountCtx)
	r.Get("/", rs.get)
	r.Put("/", rs.update)
	r.Delete("/", rs.delete)
	r.Route("/token/{tokenID}", func(r chi.Router) {
		r.Put("/", rs.updateToken)
		r.Delete("/", rs.deleteToken)
	})
	r.Put("/profile", rs.updateProfile)
	return r
}

func (rs *AccountResource) accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := auth.ClaimsFromCtx(r.Context())
		logging.GetLogEntry(r).WithField("account_id", claims.ID)
		account, err := rs.Store.Get(claims.ID)
		if err != nil {
			// account deleted while access token still valid
			logging.GetLogEntry(r).WithField("account", claims.Sub).Warn(err)
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxAccount, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type accountRequest struct {
	*auth.Account
	// not really neccessary here as we limit updated database columns in store
	ProtectedID     int      `json:"id"`
	ProtectedActive bool     `json:"active"`
	ProtectedRoles  []string `json:"roles"`
}

func (d *accountRequest) Bind(r *http.Request) error {
	// d.ProtectedActive = true
	// d.ProtectedRoles = []string{}
	return nil
}

type accountResponse struct {
	*auth.Account
}

func newAccountResponse(a *auth.Account) *accountResponse {
	resp := &accountResponse{Account: a}
	return resp
}

func (rs *AccountResource) get(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*auth.Account)
	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) update(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*auth.Account)
	data := &accountRequest{Account: acc}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	acc = data.Account
	if err := rs.Store.Update(acc); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(err))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newAccountResponse(acc))
}

func (rs *AccountResource) delete(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*auth.Account)
	if err := rs.Store.Delete(acc); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

type tokenRequest struct {
	Identifier  string
	ProtectedID int `json:"id"`
}

func (d *tokenRequest) Bind(r *http.Request) error {
	d.Identifier = strings.TrimSpace(d.Identifier)
	return nil
}

func (rs *AccountResource) updateToken(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tokenID"))
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	data := &tokenRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	acc := r.Context().Value(ctxAccount).(*auth.Account)
	for _, t := range acc.Token {
		if t.ID == id {
			if err := rs.Store.UpdateToken(&auth.Token{
				ID:         t.ID,
				Identifier: data.Identifier,
			}); err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		}
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *AccountResource) deleteToken(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tokenID"))
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	acc := r.Context().Value(ctxAccount).(*auth.Account)
	for _, t := range acc.Token {
		if t.ID == id {
			rs.Store.DeleteToken(&auth.Token{ID: t.ID})
		}
	}
	render.Respond(w, r, http.NoBody)
}

type profileRequest struct {
	*models.Profile
	ProtectedID int `json:"id"`
}

func (d *profileRequest) Bind(r *http.Request) error {
	// d.ProtectedActive = true
	// d.ProtectedRoles = []string{}
	return nil
}

type profileResponse struct {
	*models.Profile
}

func newProfileResponse(p *models.Profile) *profileResponse {
	return &profileResponse{
		Profile: p,
	}
}

func (rs *AccountResource) updateProfile(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxAccount).(*auth.Account)
	data := &profileRequest{Profile: acc.Profile}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	p := data.Profile
	if err := rs.Store.UpdateProfile(p); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(err))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newProfileResponse(p))
}
