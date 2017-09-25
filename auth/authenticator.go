package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

type ctxKey int

const (
	ctxClaims ctxKey = iota
	ctxRefreshToken
)

var (
	errTokenUnauthorized   = errors.New("token unauthorized")
	errTokenExpired        = errors.New("token expired")
	errInvalidAccessToken  = errors.New("invalid access token")
	errInvalidRefreshToken = errors.New("invalid refresh token")
)

// ClaimsFromCtx retrieves the parsed AppClaims from request context.
func ClaimsFromCtx(ctx context.Context) AppClaims {
	return ctx.Value(ctxClaims).(AppClaims)
}

// RefreshTokenFromCtx retrieves the parsed refresh token from context.
func RefreshTokenFromCtx(ctx context.Context) string {
	return ctx.Value(ctxRefreshToken).(string)
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			log(r).Warn(err)
			render.Render(w, r, ErrUnauthorized(errTokenUnauthorized))
			return
		}

		if !token.Valid {
			render.Render(w, r, ErrUnauthorized(errTokenExpired))
			return
		}

		// Token is authenticated, parse claims
		pc, ok := parseClaims(claims)
		if !ok {
			render.Render(w, r, ErrUnauthorized(errInvalidAccessToken))
			return
		}
		ctx := context.WithValue(r.Context(), ctxClaims, pc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateRefreshJWT checks validity of refresh tokens and is only used for access token refresh and logout requests. It responds with 401 Unauthorized for invalid or expired refresh tokens.
func AuthenticateRefreshJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			log(r).Warn(err)
			render.Render(w, r, ErrUnauthorized(errTokenUnauthorized))
			return
		}
		if !token.Valid {
			render.Render(w, r, ErrUnauthorized(errTokenExpired))
			return
		}
		refreshToken, ok := parseRefreshClaims(claims)
		if !ok {
			render.Render(w, r, ErrUnauthorized(errInvalidRefreshToken))
			return
		}
		// Token is authenticated, set on context
		ctx := context.WithValue(r.Context(), ctxRefreshToken, refreshToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
