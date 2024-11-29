package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
)

type RegisterPayload struct {
	username string
	email    string
	password string
	roleid   string
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload RegisterPayload
		if err := c.Bind(&payload); err != nil {
			log.Errorf("Error while binding payload <- %v", err)
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body", Errors: []string{err.Error()}})
		}

		if len(payload.roleid) < 1 {
			log.Debug("Role ID not provided, setting to default")
			payload.roleid = "1"
		}

		if models.UserExists(payload.email) {
			log.Errorf("User already exists")
			return c.JSON(http.StatusConflict, models.JSONErrorResponse{Code: http.StatusConflict, Message: "User already exists"})
		}

		user, err := models.CreateUser(&models.User{Id: uuid.NewV4().String(), Username: payload.username, Email: payload.email}, payload.password, payload.roleid)
		if err != nil {
			log.Errorf("Error creating user <- %v", err)
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: "Error creating user", Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusCreated, user)
	}
}

type LoginInfo struct {
	Token string       `json:"token"`
	Otp   string       `json:"otp"`
	User  *models.User `json:"user"`
}

func Login(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {

		var payload LoginPayload
		if err := c.Bind(&payload); err != nil {
			log.Errorf("Error while binding payload <- %v", err)
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body", Errors: []string{err.Error()}})
		}

		user, err := models.GetUserFromEmail(payload.Email)
		if err != nil {
			log.Errorf("Error while getting user from email <- %v", err)
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{Code: http.StatusNotFound, Message: fmt.Sprintf("User not found. Cause -> %v", err)})
		}

		err = user.VerifyPassword(payload.Password)
		if err != nil {
			log.Errorf("Error while verifying password <- %v", err)
			return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized: wrong password. Cause -> %v", err)})
		}

		token, err := user.GenerateToken()
		if err != nil {
			log.Errorf("Error while generating token <- %v", err)
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error while generating token. Cause -> %v", err)})
		}

		otp := cm.GenerateNewOtp()

		userToReturn, err := user.ToUser()
		if err != nil {
			log.Errorf("Error while converting user to return <- %v", err)
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error while converting user to return. Cause -> %v", err)})
		}

		return c.JSON(http.StatusOK, LoginInfo{Token: token, Otp: otp, User: userToReturn})
	}
}

type TokenInfo struct {
	Token string `json:"token"`
}

type CheckResponse struct {
	Valid bool   `json:"valid"`
	Otp   string `json:"otp"`
}

func CheckToken(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload TokenInfo
		if err := c.Bind(&payload); err != nil {
			log.Errorf("Error while binding payload <- %v", err)
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: "Invalid request body", Errors: []string{err.Error()}})
		}

		_, err := helpers.ValidateToken(payload.Token)
		if err != nil {
			log.Errorf("Error while validating token <- %v", err)
			return c.JSON(http.StatusUnauthorized, models.JSONErrorResponse{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Unauthorized. Cause -> %v", err)})
		}

		otp := cm.GenerateNewOtp()

		return c.JSON(http.StatusOK, CheckResponse{Valid: true, Otp: otp})
	}
}
