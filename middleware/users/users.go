package users

import (
	"app/models"
	"fmt"
	"github.com/labstack/echo"
	"strings"
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

		sub := getClaim(claims, "sub")

		if sub == "" {
			return nil
		}

		name := getClaim(claims, "first_name")
		family_name := getClaim(claims, "family_name")
		email := getClaim(claims, "email")
		avatar := getClaim(claims, "avatar")

		user := models.User{}

		// Check if the user exists with the auth0 user_id
		if models.DB.Where(&models.User{Auth0Id: sub}).First(&user).RecordNotFound() {

			id_parts := strings.SplitN(sub, "|", 2)

			provider := id_parts[0]
			provider_id := id_parts[1]

			// check if legacy user exists
			if models.DB.Where(&models.User{Provider: provider, ProviderId: provider_id}).First(&user).RecordNotFound() {

				// it does not exists, lets create a new one
				models.DB.Create(&models.User{
					Auth0Id:    sub,
					Name:       name + " " + family_name,
					Email:      email,
					Avatar:     avatar,
					Provider:   provider,
					ProviderId: provider_id,
				})

			} else {
				// it exists so lets link to the auth0 user_id
				user.Auth0Id = sub
				models.DB.Model(&user).Updates(user)
			}

		}

		c.Set("User", user)

		return nil

	}

}

func getClaim(claims map[string]interface{}, key string) string {
	if claims[key] == nil {
		return ""
	}
	return claims[key].(string)
}
