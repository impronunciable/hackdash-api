package middleware

import (
	"app/middleware/auth"
	"app/middleware/users"
	"github.com/labstack/echo"
)

// we need this since echo.Middleware is an empty interface
type ActualMiddleware interface {
	Handler() echo.HandlerFunc
}

type MultiMiddleware struct {
	ms []ActualMiddleware
}

func NewMultiMiddleware(ms ...ActualMiddleware) MultiMiddleware {
	return MultiMiddleware{ms}
}

func (m *MultiMiddleware) Handler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		var err error = nil
		for i := range m.ms {
			if err != nil {
				break
			}
			err = m.ms[i].Handler()(c)
		}
		return err
	}
}

func NewRestrictMiddleware(signingKey []byte, aud string) MultiMiddleware {
	authM := auth.NewJwtMiddleware("Bearer", signingKey, aud)
	usersM := users.NewUsersMiddleware()
	return NewMultiMiddleware(authM, usersM)
}
