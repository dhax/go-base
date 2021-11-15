package jwt

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
)

// TokenAuth implements JWT authentication flow.
type TokenAuth struct {
	JwtAuth          *jwtauth.JWTAuth
	JwtExpiry        time.Duration
	JwtRefreshExpiry time.Duration
}

// NewTokenAuth configures and returns a JWT authentication instance.
func NewTokenAuth() (*TokenAuth, error) {
	secret := viper.GetString("auth_jwt_secret")
	if secret == "random" {
		secret = randStringBytes(32)
	}

	a := &TokenAuth{
		JwtAuth:          jwtauth.New("HS256", []byte(secret), nil),
		JwtExpiry:        viper.GetDuration("auth_jwt_expiry"),
		JwtRefreshExpiry: viper.GetDuration("auth_jwt_refresh_expiry"),
	}

	return a, nil
}

// Verifier http middleware will verify a jwt string from a http request.
func (a *TokenAuth) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(a.JwtAuth)
}

// GenTokenPair returns both an access token and a refresh token.
func (a *TokenAuth) GenTokenPair(accessClaims AppClaims, refreshClaims RefreshClaims) (string, string, error) {
	access, err := a.CreateJWT(accessClaims)
	if err != nil {
		return "", "", err
	}
	refresh, err := a.CreateRefreshJWT(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// CreateJWT returns an access token for provided account claims.
func (a *TokenAuth) CreateJWT(c AppClaims) (string, error) {
	c.IssuedAt = time.Now().Unix()
	c.ExpiresAt = time.Now().Add(a.JwtExpiry).Unix()

	claims, err := ParseStructToMap(c)
	if err != nil {
		return "", err
	}

	_, tokenString, err := a.JwtAuth.Encode(claims)
	return tokenString, err
}

func ParseStructToMap(c interface{}) (map[string]interface{}, error) {
	var claims map[string]interface{}
	inrec, _ := json.Marshal(c)
	err := json.Unmarshal(inrec, &claims)
	if err != nil {
		return nil, err
	}

	return claims, err
}

// CreateRefreshJWT returns a refresh token for provided token Claims.
func (a *TokenAuth) CreateRefreshJWT(c RefreshClaims) (string, error) {
	c.IssuedAt = time.Now().Unix()
	c.ExpiresAt = time.Now().Add(a.JwtRefreshExpiry).Unix()

	claims, err := ParseStructToMap(c)
	if err != nil {
		return "", err
	}

	_, tokenString, err := a.JwtAuth.Encode(claims)
	return tokenString, err
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
