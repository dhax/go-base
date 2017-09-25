package auth

import (
	"net/http"
	"time"

	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/dhax/go-base/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Store defines database operations on account and token data.
type Store interface {
	GetByID(id int) (*models.Account, error)
	GetByEmail(email string) (*models.Account, error)
	GetByRefreshToken(token string) (*models.Account, *models.Token, error)
	UpdateAccount(a *models.Account) error
	SaveRefreshToken(u *models.Token) error
	DeleteRefreshToken(t *models.Token) error
	PurgeExpiredToken() error
}

// EmailService defines methods to send account emails.
type EmailService interface {
	LoginToken(name, email string, c email.LoginTokenContent) error
}

// Resource implements passwordless token authentication against a database.
type Resource struct {
	Login  *LoginTokenAuth
	Token  *TokenAuth
	store  Store
	mailer EmailService
}

// NewResource returns a configured authentication resource.
func NewResource(store Store, mailer EmailService) (*Resource, error) {
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

	resource.Cleanup()

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

func (rs *Resource) Cleanup() {
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
