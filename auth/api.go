package auth

import (
	"net/http"
	"time"

	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Storer defines database operations on account and token data.
type Storer interface {
	GetByID(id int) (*Account, error)
	GetByEmail(email string) (*Account, error)
	GetByRefreshToken(token string) (*Account, *Token, error)
	UpdateAccount(a *Account) error
	SaveRefreshToken(t *Token) error
	DeleteRefreshToken(t *Token) error
	PurgeExpiredToken() error
}

// Mailer defines methods to send account emails.
type Mailer interface {
	LoginToken(name, email string, c email.ContentLoginToken) error
}

// Resource implements passwordless token authentication against a database.
type Resource struct {
	Login  *LoginTokenAuth
	Token  *TokenAuth
	store  Storer
	mailer Mailer
}

// NewResource returns a configured authentication resource.
func NewResource(store Storer, mailer Mailer) (*Resource, error) {
	loginAuth, err := NewLoginTokenAuth()
	if err != nil {
		return nil, err
	}

	tokenAuth, err := NewTokenAuth()
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		Login:  loginAuth,
		Token:  tokenAuth,
		store:  store,
		mailer: mailer,
	}

	resource.cleanupTicker()

	return resource, nil
}

// Router provides neccessary routes for passwordless authentication flow.
func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/login", rs.login)
	r.Post("/token", rs.token)
	r.Group(func(r chi.Router) {
		r.Use(rs.Token.Verifier())
		r.Use(AuthenticateRefreshJWT)
		r.Post("/refresh", rs.refresh)
		r.Post("/logout", rs.logout)
	})
	return r
}

func (rs *Resource) cleanupTicker() {
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if err := rs.store.PurgeExpiredToken(); err != nil {
				logging.Logger.WithField("auth", "cleanup").Error(err)
			}
		}
	}()
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
