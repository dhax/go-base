package auth

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/mssola/user_agent"
	uuid "github.com/satori/go.uuid"

	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/models"
)

// The list of error types presented to the end user as error message.
var (
	ErrInvalidLogin  = errors.New("invalid email address")
	ErrUnknownLogin  = errors.New("email not registered")
	ErrLoginDisabled = errors.New("login for account disabled")
	ErrLoginToken    = errors.New("invalid or expired login token")
)

type loginRequest struct {
	Email string
}

func (body *loginRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	if err := validation.ValidateStruct(body,
		validation.Field(&body.Email, validation.Required, is.Email),
	); err != nil {
		return err
	}
	return nil
}

func (rs *Resource) login(w http.ResponseWriter, r *http.Request) {
	body := &loginRequest{}
	if err := render.Bind(r, body); err != nil {
		log(r).WithField("email", body.Email).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	acc, err := rs.store.GetByEmail(body.Email)
	if err != nil {
		log(r).WithField("email", body.Email).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	lt := rs.Login.CreateToken(acc.ID)

	go func() {
		err := rs.mailer.LoginToken(acc.Name, acc.Email, email.LoginTokenContent{
			Email:  acc.Email,
			Name:   acc.Name,
			URL:    path.Join(rs.Login.loginURL, lt.Token),
			Token:  lt.Token,
			Expiry: lt.Expiry,
		})
		if err != nil {
			log(r).WithField("module", "email").Error(err.Error())
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

	if err := validation.ValidateStruct(body,
		validation.Field(&body.Token, validation.Required, is.Alphanumeric),
	); err != nil {
		return err
	}
	return nil
}

func (rs *Resource) token(w http.ResponseWriter, r *http.Request) {
	body := &tokenRequest{}
	if err := render.Bind(r, body); err != nil {
		log(r).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrLoginToken))
		return
	}

	id, err := rs.Login.GetAccountID(body.Token)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrLoginToken))
		return
	}

	acc, err := rs.store.GetByID(id)
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
	token := &models.Token{
		Token:      uuid.NewV4().String(),
		Expiry:     time.Now().Add(time.Minute * rs.Token.jwtRefreshExpiry),
		UpdatedAt:  time.Now(),
		AccountID:  acc.ID,
		Mobile:     ua.Mobile(),
		Identifier: fmt.Sprintf("%s on %s", browser, ua.OS()),
	}

	if err := rs.store.SaveRefreshToken(token); err != nil {
		log(r).Error(err)
		render.Respond(w, r, ErrInternalServerError)
		return
	}

	access, refresh := rs.Token.GenTokenPair(acc, token)

	acc.LastLogin = time.Now()
	if err := rs.store.UpdateAccount(acc); err != nil {
		log(r).Error(err)
		render.Respond(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
}

func (rs *Resource) refresh(w http.ResponseWriter, r *http.Request) {
	rt := RefreshTokenFromCtx(r.Context())

	acc, token, err := rs.store.GetByRefreshToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(errTokenExpired))
		return
	}

	if time.Now().After(token.Expiry) {
		rs.store.DeleteRefreshToken(token)
		render.Render(w, r, ErrUnauthorized(errTokenExpired))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	token.Token = uuid.NewV4().String()
	token.Expiry = time.Now().Add(time.Minute * rs.Token.jwtRefreshExpiry)
	token.UpdatedAt = time.Now()

	access, refresh := rs.Token.GenTokenPair(acc, token)
	if err := rs.store.SaveRefreshToken(token); err != nil {
		log(r).Error(err)
		render.Respond(w, r, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.store.UpdateAccount(acc); err != nil {
		log(r).Error(err)
		render.Respond(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
}

func (rs *Resource) logout(w http.ResponseWriter, r *http.Request) {
	rt := RefreshTokenFromCtx(r.Context())
	_, token, err := rs.store.GetByRefreshToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(errTokenExpired))
		return
	}
	rs.store.DeleteRefreshToken(token)

	render.Respond(w, r, http.NoBody)
}
