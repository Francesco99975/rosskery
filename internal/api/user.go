package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
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

func User(id string) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := models.GetUserById(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("User not found. Cause -> %v", err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

func UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req models.User
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		user, err := models.GetUserById(req.Id)
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

		if err := user.Delete(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error deleting user: %v", err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}


