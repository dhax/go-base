package authorize

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/dhax/go-base/auth/jwt"
)

// RequiresRole middleware restricts access to accounts having role parameter in their jwt claims.
func RequiresRole(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			claims := jwt.ClaimsFromCtx(r.Context())
			if !hasRole(role, claims.Roles) {
				render.Render(w, r, ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func hasRole(role string, roles []string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
