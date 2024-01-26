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
	Email string `json:"email"`
	Password string `json:"password"`
}

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload RegisterPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body", Errors: []string{err.Error()}})
		}

		if len(payload.roleid) < 1 {
			payload.roleid = "1"
		}

		if models.UserExists(payload.email) {
			return c.JSON(http.StatusConflict, models.JSONErrorResponse{ Code: http.StatusConflict, Message: "User already exists"})
		}

		user, err := models.CreateUser(&models.User{Id: uuid.NewV4().String(), Username: payload.username, Email: payload.email}, payload.password, payload.roleid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{ Code: http.StatusInternalServerError, Message: "Error creating user", Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusCreated, user)
	}
}


type LoginInfo struct{
	Token string `json:"token"`
	Otp string	`json:"otp"`
}


func Login(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {

		var payload LoginPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body", Errors: []string{err.Error()}})
		}

		user, err := models.GetUserFromEmail(payload.Email)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{ Code: http.StatusNotFound, Message: fmt.Sprintf("User not found. Cause -> %v", err)})
		}

		err = user.VerifyPassword(payload.Password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{ Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized: wrong password. Cause -> %v", err)})
		}

		token, err := user.GenerateToken()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{ Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error while generating token. Cause -> %v", err)})
		}

		otp := cm.GenerateNewOtp()

		return c.JSON(http.StatusOK, LoginInfo{ Token: token, Otp: otp})
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
