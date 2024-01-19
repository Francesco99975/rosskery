package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"

	"github.com/Francesco99975/rosskery/internal/models"
)

type RegisterPayload struct {
	username string
	email string
	password string
	roleid string
}

type LoginPayload struct {
	email string
	password string
}

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload RegisterPayload
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		if len(payload.roleid) < 1 {
			payload.roleid = "1"
		}

		if models.UserExists(payload.email) {
			return echo.NewHTTPError(http.StatusConflict, "User already exists")
		}

		user, err := models.CreateUser(&models.User{ Id: uuid.NewV4().String(), Username: payload.username, Email: payload.email}, payload.password, payload.roleid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, user)
	}
}

func Login(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {

		var payload LoginPayload
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error binding payload: %v", err))
		}

		user, err := models.GetUserFromEmail(payload.email)

		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User not found. Cause -> %v", err))
		}

		err = user.VerifyPassword(payload.password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized: wrong password. Cause -> %v", err))
		}

		token, err := user.GenerateToken()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while generating token. Cause -> %v", err))
		}

		otp := cm.GenerateNewOtp()

		return c.JSON(http.StatusOK, struct{ token string; otp string }{token: token, otp: otp})
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
