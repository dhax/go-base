// Package pwdless provides JSON Web Token (JWT) authentication and authorization middleware.
// It implements a passwordless authentication flow by sending login tokens vie email which are then exchanged for JWT access and refresh tokens.
package pwdless

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/dhax/go-base/auth/jwt"
	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofrs/uuid"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
)

// AuthStorer defines database operations on accounts and tokens.
type AuthStorer interface {
	GetAccount(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	UpdateAccount(a *Account) error

	GetToken(token string) (*jwt.Token, error)
	CreateOrUpdateToken(t *jwt.Token) error
	DeleteToken(t *jwt.Token) error
	PurgeExpiredToken() error
}

// Mailer defines methods to send account emails.
type Mailer interface {
	LoginToken(name, email string, c email.ContentLoginToken) error
}

// Resource implements passwordless account authentication against a database.
type Resource struct {
	LoginAuth *LoginTokenAuth
	TokenAuth *jwt.TokenAuth
	Store     AuthStorer
	Mailer    Mailer
}

// NewResource returns a configured authentication resource.
func NewResource(authStore AuthStorer, mailer Mailer) (*Resource, error) {
	loginAuth, err := NewLoginTokenAuth()
	if err != nil {
		return nil, err
	}

	tokenAuth, err := jwt.NewTokenAuth()
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		LoginAuth: loginAuth,
		TokenAuth: tokenAuth,
		Store:     authStore,
		Mailer:    mailer,
	}

	resource.choresTicker()

	return resource, nil
}

// Router provides necessary routes for passwordless authentication flow.
func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/login", rs.login)
	r.Post("/token", rs.token)
	r.Group(func(r chi.Router) {
		r.Use(rs.TokenAuth.Verifier())
		r.Use(jwt.AuthenticateRefreshJWT)
		r.Post("/refresh", rs.refresh)
		r.Post("/logout", rs.logout)
	})
	return r
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}

type loginRequest struct {
	Email string
}

func (body *loginRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	return validation.ValidateStruct(body,
		validation.Field(&body.Email, validation.Required, is.Email),
	)
}

func (rs *Resource) login(w http.ResponseWriter, r *http.Request) {
	body := &loginRequest{}
	if err := render.Bind(r, body); err != nil {
		log(r).WithField("email", body.Email).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	acc, err := rs.Store.GetAccountByEmail(body.Email)
	if err != nil {
		log(r).WithField("email", body.Email).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	lt := rs.LoginAuth.CreateToken(acc.ID)

	go func() {
		content := email.ContentLoginToken{
			Email:  acc.Email,
			Name:   acc.Name,
			URL:    path.Join(rs.LoginAuth.loginURL, lt.Token),
			Token:  lt.Token,
			Expiry: lt.Expiry,
		}
		if err := rs.Mailer.LoginToken(acc.Name, acc.Email, content); err != nil {
			log(r).WithField("module", "email").Error(err)
		}
	}()

	render.Respond(w, r, http.NoBody)
}

type tokenRequest struct {
	Token string `json:"token"`
}

type tokenResponse struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (body *tokenRequest) Bind(r *http.Request) error {
	body.Token = strings.TrimSpace(body.Token)

	return validation.ValidateStruct(body,
		validation.Field(&body.Token, validation.Required, is.Alphanumeric),
	)
}

func (rs *Resource) token(w http.ResponseWriter, r *http.Request) {
	body := &tokenRequest{}
	if err := render.Bind(r, body); err != nil {
		log(r).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrLoginToken))
		return
	}

	id, err := rs.LoginAuth.GetAccountID(body.Token)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrLoginToken))
		return
	}

	acc, err := rs.Store.GetAccount(id)
	if err != nil {
		// account deleted before login token expired
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	token := &jwt.Token{
		Token:      uuid.Must(uuid.NewV4()).String(),
		Expiry:     time.Now().Add(rs.TokenAuth.JwtRefreshExpiry),
		UpdatedAt:  time.Now(),
		AccountID:  acc.ID,
		Mobile:     ua.Mobile(),
		Identifier: fmt.Sprintf("%s on %s", browser, ua.OS()),
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.Store.UpdateAccount(acc); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
}

func (rs *Resource) refresh(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())

	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	if time.Now().After(token.Expiry) {
		rs.Store.DeleteToken(token)
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	acc, err := rs.Store.GetAccount(token.AccountID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	token.Token = uuid.Must(uuid.NewV4()).String()
	token.Expiry = time.Now().Add(rs.TokenAuth.JwtRefreshExpiry)
	token.UpdatedAt = time.Now()

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.Store.UpdateAccount(acc); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
}

func (rs *Resource) logout(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())
	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}
	rs.Store.DeleteToken(token)

	render.Respond(w, r, http.NoBody)
}
