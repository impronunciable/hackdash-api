package auth

import (
  "fmt"
  "net/http"
  "encoding/base64"
  "github.com/dgrijalva/jwt-go"
  "github.com/labstack/echo"
)

func DecodeBase64Secret(key string) ([]byte, error) {
  decoded, err := base64.URLEncoding.DecodeString(key)
  if err != nil {
    return nil, err
  }
  return decoded, nil
}

func NewJwtMiddleware(Bearer string, SigningKey []byte, Aud string) *JWTMiddleware {
  return &JWTMiddleware{
    Bearer:  Bearer,
    SigningKey: SigningKey,
    Aud: Aud,
  }
}

type JWTMiddleware struct {
  Bearer string
  Aud string
  SigningKey []byte
}

func (me *JWTMiddleware) Handler() echo.HandlerFunc { 

  return func(c *echo.Context) error {

    // Skip WebSocket
    if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
      return nil
    }

    auth := c.Request().Header.Get("Authorization")
    l := len(me.Bearer)
    he := echo.NewHTTPError(http.StatusUnauthorized)

    if len(auth) > l+1 && auth[:l] == me.Bearer {
      t, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {

        // Always check the signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
          return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        // Return the key for validation
        return me.SigningKey, nil
      })
      if err == nil && t.Valid {
        // Store token claims in echo.Context
        c.Set("claims", t.Claims)
        return nil
      }
    }
    return he
  }
}
