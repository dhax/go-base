package jwt

import (
	"errors"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

type CommonClaims struct {
	ExpiresAt int64 `json:"exp,omitempty"`
	IssuedAt  int64 `json:"iat,omitempty"`
}

// AppClaims represent the claims parsed from JWT access token.
type AppClaims struct {
	ID    int      `json:"id,omitempty"`
	Sub   string   `json:"sub,omitempty"`
	Roles []string `json:"roles,omitempty"`
	CommonClaims
}

// ParseClaims parses JWT claims into AppClaims.
func (c *AppClaims) ParseClaims(claims map[string]interface{}) error {
	id, ok := claims["id"]
	if !ok {
		return errors.New("could not parse claim id")
	}
	c.ID = int(id.(float64))

	sub, ok := claims[jwt.SubjectKey]
	if !ok {
		return errors.New("could not parse claim sub")
	}
	c.Sub = sub.(string)

	rl, ok := claims["roles"]
	if !ok {
		return errors.New("could not parse claims roles")
	}

	var roles []string
	if rl != nil {
		for _, v := range rl.([]interface{}) {
			roles = append(roles, v.(string))
		}
	}
	c.Roles = roles

	return nil
}

// RefreshClaims represents the claims parsed from JWT refresh token.
type RefreshClaims struct {
	ID    int    `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	CommonClaims
}

// ParseClaims parses the JWT claims into RefreshClaims.
func (c *RefreshClaims) ParseClaims(claims map[string]interface{}) error {
	token, ok := claims["token"]
	if !ok {
		return errors.New("could not parse claim token")
	}
	c.Token = token.(string)
	return nil
}
