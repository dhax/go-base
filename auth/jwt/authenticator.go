package jwt

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/dhax/go-base/logging"
)

type ctxKey int

const (
	ctxClaims ctxKey = iota
	ctxRefreshToken
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
			logging.GetLogEntry(r).Warn(err)
			render.Render(w, r, ErrUnauthorized(ErrTokenUnauthorized))
			return
		}

		if !token.Valid {
			render.Render(w, r, ErrUnauthorized(ErrTokenExpired))
			return
		}

		// Token is authenticated, parse claims
		var c AppClaims
		err = c.ParseClaims(claims)
		if err != nil {
			logging.GetLogEntry(r).Error(err)
			render.Render(w, r, ErrUnauthorized(ErrInvalidAccessToken))
			return
		}

		// Set AppClaims on context
		ctx := context.WithValue(r.Context(), ctxClaims, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateRefreshJWT checks validity of refresh tokens and is only used for access token refresh and logout requests. It responds with 401 Unauthorized for invalid or expired refresh tokens.
func AuthenticateRefreshJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			logging.GetLogEntry(r).Warn(err)
			render.Render(w, r, ErrUnauthorized(ErrTokenUnauthorized))
			return
		}
		if !token.Valid {
			render.Render(w, r, ErrUnauthorized(ErrTokenExpired))
			return
		}

		// Token is authenticated, parse refresh token string
		var c RefreshClaims
		err = c.ParseClaims(claims)
		if err != nil {
			logging.GetLogEntry(r).Error(err)
			render.Render(w, r, ErrUnauthorized(ErrInvalidRefreshToken))
			return
		}
		// Set refresh token string on context
		ctx := context.WithValue(r.Context(), ctxRefreshToken, c.Token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
