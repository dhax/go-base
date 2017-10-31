package pwdless

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	errTokenNotFound = errors.New("login token not found")
)

// LoginToken is an in-memory saved token referencing an account ID and an expiry date.
type LoginToken struct {
	Token     string
	AccountID int
	Expiry    time.Time
}

// LoginTokenAuth implements passwordless login authentication flow using temporary in-memory stored tokens.
type LoginTokenAuth struct {
	token            map[string]LoginToken
	mux              sync.RWMutex
	loginURL         string
	loginTokenLength int
	loginTokenExpiry time.Duration
}

// NewLoginTokenAuth configures and returns a LoginToken authentication instance.
func NewLoginTokenAuth() (*LoginTokenAuth, error) {
	a := &LoginTokenAuth{
		token:            make(map[string]LoginToken),
		loginURL:         viper.GetString("auth_login_url"),
		loginTokenLength: viper.GetInt("auth_login_token_length"),
		loginTokenExpiry: viper.GetDuration("auth_login_token_expiry"),
	}
	return a, nil
}

// CreateToken creates an in-memory login token referencing account ID. It returns a token containing a random tokenstring and expiry date.
func (a *LoginTokenAuth) CreateToken(id int) LoginToken {
	lt := LoginToken{
		Token:     randStringBytes(a.loginTokenLength),
		AccountID: id,
		Expiry:    time.Now().Add(a.loginTokenExpiry),
	}
	a.add(lt)
	a.purgeExpired()
	return lt
}

// GetAccountID looks up the token by tokenstring and returns the account ID or error if token not found or expired.
func (a *LoginTokenAuth) GetAccountID(token string) (int, error) {
	lt, exists := a.get(token)
	if !exists || time.Now().After(lt.Expiry) {
		return 0, errTokenNotFound
	}
	a.delete(lt.Token)
	return lt.AccountID, nil
}

func (a *LoginTokenAuth) get(token string) (LoginToken, bool) {
	a.mux.RLock()
	lt, ok := a.token[token]
	a.mux.RUnlock()
	return lt, ok
}

func (a *LoginTokenAuth) add(lt LoginToken) {
	a.mux.Lock()
	a.token[lt.Token] = lt
	a.mux.Unlock()
}

func (a *LoginTokenAuth) delete(token string) {
	a.mux.Lock()
	delete(a.token, token)
	a.mux.Unlock()
}

func (a *LoginTokenAuth) purgeExpired() {
	for t, v := range a.token {
		if time.Now().After(v.Expiry) {
			a.delete(t)
		}
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	for k, v := range buf {
		buf[k] = letterBytes[v%byte(len(letterBytes))]
	}
	return string(buf)
}
