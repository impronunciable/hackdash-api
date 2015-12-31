package middleware

import (
	"app/middleware/auth"
	"app/middleware/users"

	"github.com/labstack/echo"
	"net/http"
)

var authM *auth.JWTMiddleware
var usersM *users.UsersMiddleware

func middlewareAuth(c *echo.Context, forceLoggedUser bool) error {

	err := authM.Handler()(c)
	if err != nil {
		if forceLoggedUser {
			return err
		}
		return nil
	}

	if forceLoggedUser {
		he := echo.NewHTTPError(http.StatusUnauthorized)

		if c.Get("claims") == nil {
			return he
		}

		claims := c.Get("claims").(map[string]interface{})

		if _, ok := claims["sub"]; !ok {
			return he
		}
	}

	err = usersM.Handler()(c)
	return err
}

func ConfigRestrictMiddleware(signingKey []byte, aud string) {
	authM = auth.NewJwtMiddleware("Bearer", signingKey, aud)
	usersM = users.NewUsersMiddleware()
}

func Secure(delegate echo.HandlerFunc) echo.HandlerFunc {
	return middlewareWrapper(delegate, true)
}
func Unsecure(delegate echo.HandlerFunc) echo.HandlerFunc {
	return middlewareWrapper(delegate, false)
}
func middlewareWrapper(delegate echo.HandlerFunc, forceLoggedUser bool) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if err := middlewareAuth(c, forceLoggedUser); err != nil {
			return err
		} else {
			return delegate(c)
		}
	}
}
