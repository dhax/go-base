package pwdless

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// The list of error types presented to the end user as error message.
var (
	ErrInvalidLogin  = errors.New("invalid email address")
	ErrUnknownLogin  = errors.New("email not registered")
	ErrLoginDisabled = errors.New("login for account disabled")
	ErrLoginToken    = errors.New("invalid or expired login token")
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render sets the application-specific error code in AppCode.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrUnauthorized renders status 401 Unauthorized with custom error message.
func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     http.StatusText(http.StatusUnauthorized),
		ErrorText:      err.Error(),
	}
}

// The list of default error types without specific error message.
var (
	ErrInternalServerError = &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
	}
)
