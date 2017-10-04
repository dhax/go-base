package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

// AppClaims represent the claims extracted from JWT token.
type AppClaims struct {
	ID    int
	Sub   string
	Roles []string
}

// TokenAuth implements JWT authentication flow.
type TokenAuth struct {
	JwtAuth          *jwtauth.JwtAuth
	jwtExpiry        time.Duration
	jwtRefreshExpiry time.Duration
}

// NewTokenAuth configures and returns a JWT authentication instance.
func NewTokenAuth() (*TokenAuth, error) {
	secret := viper.GetString("auth_jwt_secret")
	if secret == "random" {
		secret = randStringBytes(32)
	}

	a := &TokenAuth{
		JwtAuth:          jwtauth.New("HS256", []byte(secret), nil),
		jwtExpiry:        viper.GetDuration("auth_jwt_expiry"),
		jwtRefreshExpiry: viper.GetDuration("auth_jwt_refresh_expiry"),
	}

	return a, nil
}

// Verifier http middleware will verify a jwt string from a http request.
func (a *TokenAuth) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(a.JwtAuth)
}

// GenTokenPair returns both an access token and a refresh token.
func (a *TokenAuth) GenTokenPair(ca jwtauth.Claims, cr jwtauth.Claims) (string, string, error) {
	access, err := a.CreateJWT(ca)
	if err != nil {
		return "", "", err
	}
	refresh, err := a.CreateRefreshJWT(cr)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// CreateJWT returns an access token for provided account claims.
func (a *TokenAuth) CreateJWT(c jwtauth.Claims) (string, error) {
	c.SetIssuedNow()
	c.SetExpiryIn(a.jwtExpiry * time.Minute)
	_, tokenString, err := a.JwtAuth.Encode(c)
	return tokenString, err
}

// CreateRefreshJWT returns a refresh token for provided token Claims.
func (a *TokenAuth) CreateRefreshJWT(c jwtauth.Claims) (string, error) {
	c.SetIssuedNow()
	c.SetExpiryIn(time.Minute * a.jwtRefreshExpiry)
	_, tokenString, err := a.JwtAuth.Encode(c)
	return tokenString, err
}

func parseClaims(c jwtauth.Claims) (AppClaims, bool) {
	var claims AppClaims
	allOK := true
	id, ok := c.Get("id")
	if !ok {
		allOK = false
	}
	claims.ID = int(id.(float64))

	sub, ok := c.Get("sub")
	if !ok {
		allOK = false
	}
	claims.Sub = sub.(string)

	rl, ok := c.Get("roles")
	if !ok {
		allOK = false
	}

	var roles []string
	if rl != nil {
		for _, v := range rl.([]interface{}) {
			roles = append(roles, v.(string))
		}
	}
	claims.Roles = roles

	return claims, allOK
}

func parseRefreshClaims(c jwtauth.Claims) (string, bool) {
	token, ok := c.Get("token")
	if !ok {
		return "", false
	}
	return token.(string), ok
}
