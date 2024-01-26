package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Users() echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := models.GetAllUsers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{ Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching users: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, users)
	}
}

func User() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		user, err := models.GetUserById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{ Code: http.StatusNotFound, Message: fmt.Sprintf("User not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.User
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing request body for user: %v", err), Errors: []string{err.Error()}})
		}

		user, err := models.GetUserById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{ Code: http.StatusNotFound, Message: fmt.Sprintf("User not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		if err := user.Update(&req); err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{ Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error updating user: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		user, err := models.GetUserById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{ Code: http.StatusNotFound, Message: fmt.Sprintf("User not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		defer func () {
			err = user.Delete()
			if err != nil {
				log.Errorf("Error while deleting order: %v", err)
			}
		}()

		juser, err := user.ToUser()
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{ Code: http.StatusNotFound, Message: fmt.Sprintf("User could not be converted. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, juser)
	}
}


