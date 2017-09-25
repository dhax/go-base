package jwtauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

var (
	ErrUnauthorized = errors.New("jwtauth: token is unauthorized")
	ErrExpired      = errors.New("jwtauth: token is expired")
)

type JwtAuth struct {
	signKey   []byte
	verifyKey []byte
	signer    jwt.SigningMethod
	parser    *jwt.Parser
}

// New creates a JwtAuth authenticator instance that provides middleware handlers
// and encoding/decoding functions for JWT signing.
func New(alg string, signKey []byte, verifyKey []byte) *JwtAuth {
	return NewWithParser(alg, &jwt.Parser{}, signKey, verifyKey)
}

// NewWithParser is the same as New, except it supports custom parser settings
// introduced in jwt-go/v2.4.0.
//
// We explicitly toggle `SkipClaimsValidation` in the `jwt-go` parser so that
// we can control when the claims are validated - in our case, by the Verifier
// http middleware handler.
func NewWithParser(alg string, parser *jwt.Parser, signKey []byte, verifyKey []byte) *JwtAuth {
	parser.SkipClaimsValidation = true
	return &JwtAuth{
		signKey:   signKey,
		verifyKey: verifyKey,
		signer:    jwt.GetSigningMethod(alg),
		parser:    parser,
	}
}

// Verifier http middleware handler will verify a JWT string from a http request.
//
// Verifier will search for a JWT token in a http request, in the order:
//   1. 'jwt' URI query parameter
//   2. 'Authorization: BEARER T' request header
//   3. Cookie 'jwt' value
//
// The first JWT string that is found as a query parameter, authorization header
// or cookie header is then decoded by the `jwt-go` library and a *jwt.Token
// object is set on the request context. In the case of a signature decoding error
// the Verifier will also set the error on the request context.
//
// The Verifier always calls the next http handler in sequence, which can either
// be the generic `jwtauth.Authenticator` middleware or your own custom handler
// which checks the request context jwt token and error to prepare a custom
// http response.
func Verifier(ja *JwtAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Verify(ja, "")(next)
	}
}

func Verify(ja *JwtAuth, paramAliases ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := VerifyRequest(ja, r, paramAliases...)
			ctx = NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func VerifyRequest(ja *JwtAuth, r *http.Request, paramAliases ...string) (*jwt.Token, error) {
	var tokenStr string
	var err error

	// Get token from query params
	tokenStr = r.URL.Query().Get("jwt")

	// Get token from other param aliases
	if tokenStr == "" && paramAliases != nil && len(paramAliases) > 0 {
		for _, p := range paramAliases {
			tokenStr = r.URL.Query().Get(p)
			if tokenStr != "" {
				break
			}
		}
	}

	// Get token from authorization header
	if tokenStr == "" {
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenStr = bearer[7:]
		}
	}

	// Get token from cookie
	if tokenStr == "" {
		// TODO: paramAliases should apply to cookies too..
		cookie, err := r.Cookie("jwt")
		if err == nil {
			tokenStr = cookie.Value
		}
	}

	// TODO: what other kinds of validations should we do / error messages?

	// Verify the token
	token, err := ja.Decode(tokenStr)
	if err != nil {
		switch err.Error() {
		case "token is expired":
			err = ErrExpired
		}

		// ctx = NewContext(ctx, token, err)
		// next.ServeHTTP(w, r.WithContext(ctx))
		return token, err
	}

	if token == nil || !token.Valid || token.Method != ja.signer {
		err = ErrUnauthorized
		// ctx = NewContext(ctx, token, err)
		// next.ServeHTTP(w, r.WithContext(ctx))
		return token, err
	}

	// Check expiry via "exp" claim
	if IsExpired(token) {
		err = ErrExpired
		// ctx = NewContext(ctx, token, err)
		// next.ServeHTTP(w, r.WithContext(ctx))
		return token, err
	}

	// Valid!
	return token, nil
}

func (ja *JwtAuth) Encode(claims Claims) (t *jwt.Token, tokenString string, err error) {
	t = jwt.New(ja.signer)
	t.Claims = claims
	tokenString, err = t.SignedString(ja.signKey)
	t.Raw = tokenString
	return
}

func (ja *JwtAuth) Decode(tokenString string) (t *jwt.Token, err error) {
	// Decode the tokenString, but avoid using custom Claims via jwt-go's
	// ParseWithClaims as the jwt-go types will cause some glitches, so easier
	// to decode as MapClaims then wrap the underlying map[string]interface{}
	// to our Claims type
	t, err = ja.parser.Parse(tokenString, ja.keyFunc)
	if err != nil {
		return nil, err
	}
	return
}

func (ja *JwtAuth) keyFunc(t *jwt.Token) (interface{}, error) {
	if ja.verifyKey != nil && len(ja.verifyKey) > 0 {
		return ja.verifyKey, nil
	} else {
		return ja.signKey, nil
	}
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func NewContext(ctx context.Context, t *jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func FromContext(ctx context.Context) (*jwt.Token, Claims, error) {
	token, _ := ctx.Value(TokenCtxKey).(*jwt.Token)

	var claims Claims
	if token != nil {
		tokenClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			panic("jwtauth: expecting jwt.MapClaims")
		}
		claims = Claims(tokenClaims)
	} else {
		claims = Claims{}
	}

	err, _ := ctx.Value(ErrorCtxKey).(error)

	return token, claims, err
}

func IsExpired(t *jwt.Token) bool {
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		panic("jwtauth: expecting jwt.MapClaims")
	}

	if expv, ok := claims["exp"]; ok {
		var exp int64
		switch v := expv.(type) {
		case float64:
			exp = int64(v)
		case int64:
			exp = v
		case json.Number:
			exp, _ = v.Int64()
		default:
		}

		if exp < EpochNow() {
			return true
		}
	}

	return false
}

// Claims is a convenience type to manage a JWT claims hash.
type Claims map[string]interface{}

// NOTE: as of v3.0 of jwt-go, Valid() interface method is called to verify
// the claims. However, the current design we test these claims in the
// Verifier middleware, so we skip this step.
func (c Claims) Valid() error {
	return nil
}

func (c Claims) Set(k string, v interface{}) Claims {
	c[k] = v
	return c
}

func (c Claims) Get(k string) (interface{}, bool) {
	v, ok := c[k]
	return v, ok
}

// Set issued at ("iat") to specified time in the claims
func (c Claims) SetIssuedAt(tm time.Time) Claims {
	c["iat"] = tm.UTC().Unix()
	return c
}

// Set issued at ("iat") to present time in the claims
func (c Claims) SetIssuedNow() Claims {
	c["iat"] = EpochNow()
	return c
}

// Set expiry ("exp") in the claims and return itself so it can be chained
func (c Claims) SetExpiry(tm time.Time) Claims {
	c["exp"] = tm.UTC().Unix()
	return c
}

// Set expiry ("exp") in the claims to some duration from the present time
// and return itself so it can be chained
func (c Claims) SetExpiryIn(tm time.Duration) Claims {
	c["exp"] = ExpireIn(tm)
	return c
}

// Helper function that returns the NumericDate time value used by the spec
func EpochNow() int64 {
	return time.Now().UTC().Unix()
}

// Helper function to return calculated time in the future for "exp" claim.
func ExpireIn(tm time.Duration) int64 {
	return EpochNow() + int64(tm.Seconds())
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}
