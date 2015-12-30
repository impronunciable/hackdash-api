package middleware

import (
	"app/middleware/auth"
	"app/middleware/users"

	"github.com/labstack/echo"
)

var authM *auth.JWTMiddleware
var usersM *users.UsersMiddleware

func middlewareAuth(c *echo.Context) error {
	err := authM.Handler()(c)
	if err != nil {
		return err
	}
	err = usersM.Handler()(c)
	return err
}

func ConfigRestrictMiddleware(signingKey []byte, aud string) {
	authM = auth.NewJwtMiddleware("Bearer", signingKey, aud)
	usersM = users.NewUsersMiddleware()
}

func Secure(delegate echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if err := middlewareAuth(c); err != nil {
			return err
		} else {
			return delegate(c)
		}
	}

}
