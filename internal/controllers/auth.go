package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"

	"github.com/Francesco99975/rosskery/internal/models"
)

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")
		roleid := c.FormValue("role")

		if len(roleid) < 1 {
			roleid = "1"
		}

		if models.UserExists(email) {
			return echo.NewHTTPError(http.StatusConflict, "User already exists")
		}

		user, err := models.CreateUser(&models.User{ Id: uuid.NewV4().String(), Username: username, Email: email}, password, roleid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, user)
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		email := c.FormValue("email")
		password := c.FormValue("password")

		user, err := models.GetUserFromEmail(email)

		if err != nil  {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User not found. Cause -> %v", err))
		}

		err = user.VerifyPassword(password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized: wrong password. Cause -> %v", err))
		}

		token, err := user.GenerateToken()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while generating token. Cause -> %v", err))
		}

		return c.JSON(http.StatusOK, struct { token string }{ token: token })
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
