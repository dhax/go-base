package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dhax/go-base/email"
	"github.com/dhax/go-base/logging"
	"github.com/dhax/go-base/models"
	"github.com/dhax/go-base/testing/mock"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var (
	auth      *Resource
	authstore mock.AuthStore
	mailer    mock.EmailService
	ts        *httptest.Server
)

func TestMain(m *testing.M) {
	viper.SetDefault("auth_login_token_length", 8)
	viper.SetDefault("auth_login_token_expiry", 11)
	viper.SetDefault("auth_jwt_secret", "random")
	viper.SetDefault("log_level", "error")

	var err error
	auth, err = NewResource(&authstore, &mailer)
	if err != nil {
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(logging.NewStructuredLogger(logging.NewLogger()))
	r.Mount("/", auth.Router())

	ts = httptest.NewServer(r)

	code := m.Run()
	ts.Close()
	os.Exit(code)
}

func TestAuthResource_login(t *testing.T) {
	authstore.GetByEmailFn = func(email string) (*models.Account, error) {
		var err error
		a := models.Account{
			ID:    1,
			Email: email,
			Name:  "test",
		}

		switch email {
		case "not@exists.io":
			err = errors.New("sql no row")
		case "disabled@account.io":
			a.Active = false
		case "valid@account.io":
			a.Active = true
		}
		return &a, err
	}

	mailer.LoginTokenFn = func(n, e string, c email.LoginTokenContent) error {
		return nil
	}

	tests := []struct {
		name   string
		email  string
		status int
		err    error
	}{
		{"missing", "", http.StatusUnauthorized, ErrInvalidLogin},
		{"inexistent", "not@exists.io", http.StatusUnauthorized, ErrUnknownLogin},
		{"disabled", "disabled@account.io", http.StatusUnauthorized, ErrLoginDisabled},
		{"valid", "valid@account.io", http.StatusOK, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := encode(&loginRequest{Email: tc.email})
			if err != nil {
				t.Fatal("failed to encode request body")
			}
			res, body := testRequest(t, ts, "POST", "/login", req, "")

			if res.StatusCode != tc.status {
				t.Errorf("got http status %d, want: %d", res.StatusCode, tc.status)
			}
			if tc.err != nil && !strings.Contains(body, tc.err.Error()) {
				t.Errorf(" got: %s, expected to contain: %s", body, tc.err.Error())
			}
			if tc.err == ErrInvalidLogin && authstore.GetByEmailInvoked {
				t.Error("GetByLoginToken invoked for invalid email")
			}
			if tc.err == nil && !mailer.LoginTokenInvoked {
				t.Error("emailService.LoginToken not invoked")
			}
			authstore.GetByEmailInvoked = false
			mailer.LoginTokenInvoked = false
		})
	}
}

func TestAuthResource_token(t *testing.T) {
	authstore.GetByIDFn = func(id int) (*models.Account, error) {
		var err error
		a := models.Account{
			ID:     id,
			Active: true,
			Name:   "test",
		}
		switch id {
		case 2:
			a.Active = false
		case 3:
			// unmodified
		default:
			err = errors.New("sql no rows")
		}
		return &a, err
	}
	authstore.UpdateAccountFn = func(a *models.Account) error {
		a.LastLogin = time.Now()
		return nil
	}
	authstore.SaveRefreshTokenFn = func(a *models.Token) error {
		return nil
	}

	tests := []struct {
		name   string
		token  string
		id     int
		status int
		err    error
	}{
		{"invalid", "#§$%", 0, http.StatusUnauthorized, ErrLoginToken},
		{"expired", "12345678", 0, http.StatusUnauthorized, ErrLoginToken},
		{"deleted_account", "", 1, http.StatusUnauthorized, ErrUnknownLogin},
		{"disabled", "", 2, http.StatusUnauthorized, ErrLoginDisabled},
		{"valid", "", 3, http.StatusOK, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := auth.Login.CreateToken(tc.id)
			if tc.token != "" {
				token.Token = tc.token
			}

			req, err := encode(tokenRequest{Token: token.Token})
			if err != nil {
				t.Fatal("failed to encode request body")
			}
			res, body := testRequest(t, ts, "POST", "/token", req, "")

			if res.StatusCode != tc.status {
				t.Errorf("got http status %d, want: %d", res.StatusCode, tc.status)
			}
			if tc.err != nil && !strings.Contains(body, tc.err.Error()) {
				t.Errorf("got: %s, expected to contain: %s", body, tc.err.Error())
			}
			if tc.err == ErrLoginToken && authstore.SaveRefreshTokenInvoked {
				t.Errorf("SaveRefreshToken invoked despite error %s", tc.err.Error())
			}
			if tc.err == nil && !authstore.SaveRefreshTokenInvoked {
				t.Error("SaveRefreshToken not invoked")
			}
			authstore.SaveRefreshTokenInvoked = false
		})
	}
}

func TestAuthResource_refresh(t *testing.T) {
	authstore.GetByRefreshTokenFn = func(token string) (*models.Account, *models.Token, error) {
		var err error
		a := models.Account{
			Active: true,
			Name:   "Test",
		}
		var t models.Token
		t.Expiry = time.Now().Add(1 * time.Minute)

		switch token {
		case "notfound":
			err = errors.New("sql no rows")
		case "expired":
			t.Expiry = time.Now().Add(-1 * time.Minute)
		case "disabled":
			a.Active = false
		case "valid":
			// unmodified
		}
		return &a, &t, err
	}
	authstore.UpdateAccountFn = func(a *models.Account) error {
		a.LastLogin = time.Now()
		return nil
	}
	authstore.SaveRefreshTokenFn = func(a *models.Token) error {
		return nil
	}
	authstore.DeleteRefreshTokenFn = func(t *models.Token) error {
		return nil
	}

	tests := []struct {
		name   string
		token  string
		exp    time.Duration
		status int
		err    error
	}{
		{"notfound", "notfound", 1, http.StatusUnauthorized, errTokenExpired},
		{"expired", "expired", -1, http.StatusUnauthorized, errTokenUnauthorized},
		{"disabled", "disabled", 1, http.StatusUnauthorized, ErrLoginDisabled},
		{"valid", "valid", 1, http.StatusOK, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jwt := genJWT(jwtauth.Claims{"token": tc.token, "exp": time.Minute * tc.exp})
			res, body := testRequest(t, ts, "POST", "/refresh", nil, jwt)
			if res.StatusCode != tc.status {
				t.Errorf("got http status %d, want: %d", res.StatusCode, tc.status)
			}
			if tc.err != nil && !strings.Contains(body, tc.err.Error()) {
				t.Errorf("got: %s, expected error to contain: %s", body, tc.err.Error())
			}
			if tc.status == http.StatusUnauthorized && authstore.SaveRefreshTokenInvoked {
				t.Errorf("SaveRefreshToken invoked for status %d", tc.status)
			}
			if tc.status == http.StatusOK && !authstore.GetByRefreshTokenInvoked {
				t.Errorf("GetRefreshToken not invoked")
			}
			if tc.status == http.StatusOK && !authstore.SaveRefreshTokenInvoked {
				t.Errorf("SaveRefreshToken not invoked")
			}
			if tc.status == http.StatusOK && authstore.DeleteRefreshTokenInvoked {
				t.Errorf("DeleteRefreshToken should not be invoked")
			}
			authstore.GetByRefreshTokenInvoked = false
			authstore.SaveRefreshTokenInvoked = false
			authstore.DeleteRefreshTokenInvoked = false
		})
	}
}

func TestAuthResource_logout(t *testing.T) {
	authstore.GetByRefreshTokenFn = func(token string) (*models.Account, *models.Token, error) {
		var err error
		var a models.Account
		t := models.Token{
			Expiry: time.Now().Add(1 * time.Minute),
		}

		switch token {
		case "notfound":
			err = errors.New("sql no rows")
		}
		return &a, &t, err
	}
	authstore.DeleteRefreshTokenFn = func(a *models.Token) error {
		return nil
	}

	tests := []struct {
		name   string
		token  string
		exp    time.Duration
		status int
		err    error
	}{
		{"notfound", "notfound", 1, http.StatusUnauthorized, errTokenExpired},
		{"expired", "valid", -1, http.StatusUnauthorized, errTokenUnauthorized},
		{"valid", "valid", 1, http.StatusOK, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jwt := genJWT(jwtauth.Claims{"token": tc.token, "exp": time.Minute * tc.exp})
			res, body := testRequest(t, ts, "POST", "/logout", nil, jwt)
			if res.StatusCode != tc.status {
				t.Errorf("got http status %d, want: %d", res.StatusCode, tc.status)
			}
			if tc.err != nil && !strings.Contains(body, tc.err.Error()) {
				t.Errorf("got: %x, expected error to contain %s", body, tc.err.Error())
			}
			if tc.status == http.StatusUnauthorized && authstore.DeleteRefreshTokenInvoked {
				t.Errorf("DeleteRefreshToken invoked for status %d", tc.status)
			}
			authstore.DeleteRefreshTokenInvoked = false
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, token string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "BEARER "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	return resp, string(respBody)
}

func genJWT(c jwtauth.Claims) string {
	_, tokenString, _ := auth.Token.JwtAuth.Encode(c)
	return tokenString
}

func encode(v interface{}) (*bytes.Buffer, error) {
	data := new(bytes.Buffer)
	err := json.NewEncoder(data).Encode(v)
	return data, err
}
