package auth

import (
	"net/http"

	"github.com/go-chi/render"
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

// ErrRender returns status 422 Unprocessable Entity for invalid request body
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
	}
}

// ErrInvalidRequest returns status 422 Unprocessable Entity with validation errors
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
	}
}

// The list of default error types without specific error message.
var (
	ErrBadRequest = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)}

	ErrForbidden = &ErrResponse{HTTPStatusCode: http.StatusForbidden, StatusText: http.StatusText(http.StatusForbidden)}

	ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound)}

	ErrInternalServerError = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)}
)
