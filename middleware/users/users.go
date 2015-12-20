package users

import (
  "app/models"
  "github.com/labstack/echo"
)


func NewUsersMiddleware() *UsersMiddleware {
  return &UsersMiddleware{}
}

type UsersMiddleware struct {
}

func (me *UsersMiddleware) Handler() echo.HandlerFunc { 

  return func(c *echo.Context) error {
    
    // Skip WebSocket
    if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
      return nil
    }

    claims := c.Get("claims").(map[string]interface{})
    sub := claims["sub"].(string)

    user := models.User{}

    if models.DB.Where(&models.User{Auth0Id: sub}).First(&user).RecordNotFound() {

      models.DB.Create(&models.User{
        Auth0Id: sub,
      })

    }

    c.Set("User", user)

    return nil

  }

}
