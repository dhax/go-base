package auth

import (
	"net/http"
	"time"

	"github.com/dhax/go-base/models"
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

// GenTokenPair returns both an access token and a refresh token for provided account.
func (a *TokenAuth) GenTokenPair(u *models.Account, tok *models.Token) (string, string) {
	access := a.CreateJWT(u)
	refresh := a.CreateRefreshJWT(tok)
	return access, refresh
}

// CreateJWT returns an access token for provided account.
func (a *TokenAuth) CreateJWT(acc *models.Account) string {
	claims := jwtauth.Claims{
		"id":    acc.ID,
		"sub":   acc.Name,
		"roles": acc.Roles,
	}
	claims.SetIssuedNow()
	claims.SetExpiryIn(a.jwtExpiry * time.Minute)

	_, tokenString, _ := a.JwtAuth.Encode(claims)
	return tokenString
}

// CreateRefreshJWT returns a refresh token for provided account.
func (a *TokenAuth) CreateRefreshJWT(tok *models.Token) string {
	claims := jwtauth.Claims{
		"id":    tok.ID,
		"token": tok.Token,
	}
	claims.SetIssuedNow()
	claims.SetExpiryIn(time.Minute * a.jwtRefreshExpiry)

	_, tokenString, _ := a.JwtAuth.Encode(claims)
	return tokenString
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
