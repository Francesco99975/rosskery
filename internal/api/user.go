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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching users: %v", err))
		}

		return c.JSON(http.StatusOK, users)
	}
}

func User() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		user, err := models.GetUserById(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User not found. Cause -> %v", err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

func UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req models.User
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		user, err := models.GetUserById(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User not found. Cause -> %v", err))
		}

		if err := user.Update(&req); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error updating user: %v", err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		user, err := models.GetUserById(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User not found. Cause -> %v", err))
		}

		defer func () {
			err = user.Delete()
			if err != nil {
				log.Errorf("Error while deleting order: %v", err)
			}
		}()

		juser, err := user.ToUser()
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User could not be converted. Cause -> %v", err))
		}

		return c.JSON(http.StatusOK, juser)
	}
}


