package httpsvc

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fahmifan/devkit/pkg/core/auth"
	"github.com/labstack/echo/v4"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUnauthenticated = errors.New("unauthorized")
	ErrTokenInvalid    = errors.New("token invalid")
)

// ErrMissingAuthorization error
var ErrMissingAuthorization = errors.New("missing Authorization header")

func (server *Server) addUserToCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := parseTokenFromHeader(&c.Request().Header)
		if err != nil {
			return next(c)
		}

		authUser, ok := auth.ParseToken(server.jwtKey, token)
		if !ok {
			return next(c)
		}

		ctx := c.Request().Context()
		reqCtx := auth.CtxWithUser(ctx, authUser)
		c.SetRequest(c.Request().WithContext(reqCtx))

		return next(c)
	}
}

func (s *Server) authz(perms ...auth.Permission) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authUser, ok := auth.GetUserFromCtx(c.Request().Context())
			if !ok {
				return responseError(c, ErrUnauthorized)
			}

			for _, p := range perms {
				if authUser.Role.Granted(p) {
					return next(c)
				}
			}

			return responseError(c, ErrUnauthorized)
		}
	}
}

const authzHeader = "Authorization"

func parseTokenFromHeader(header *http.Header) (auth.JWTToken, error) {
	var token string

	authHeaders := strings.Split(header.Get(authzHeader), " ")
	if len(authHeaders) != 2 {
		return "", ErrTokenInvalid
	}

	if authHeaders[0] != "Bearer" {
		return "", ErrMissingAuthorization
	}

	token = strings.Trim(authHeaders[1], " ")
	if token == "" {
		return "", ErrMissingAuthorization
	}

	return auth.JWTToken(token), nil
}
