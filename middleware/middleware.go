package middleware

import (
	"github.com/labstack/echo/v4"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(echo.Context) bool

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc func(echo.Context)
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(echo.Context) bool {
	return false
}
