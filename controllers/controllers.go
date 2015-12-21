package controllers

import (
	"app/models"
	"github.com/labstack/echo"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func InitV3Routes(r *echo.Group) {
	InitDashboardRoutes(r)
	InitUserRoutes(r)
}

func Decode(c *echo.Context, val interface{}) error {

	if err := c.Bind(&val); err != nil {
		logger.Printf("error decoding %T, error: %s", val, err.Error())
		return err
	}

	if validable, ok := val.(models.Validable); ok {
		err := validable.Validate()
		if err != nil {
			logger.Printf("error validating %T: %v", val, err)
			return err
		}
		return nil
	} else {
		logger.Printf("%T is not validable", val)
		return nil
	}
}
