package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
)


func Roles() echo.HandlerFunc {
	return func(c echo.Context) error {
		roles, err := models.GetRoles()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{ Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching roles: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, roles)
	}
}

func Role() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Role
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing request body for role: %v", err), Errors: []string{err.Error()}})
		}

		role, err := models.GetRoleById(payload.Id)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{ Code: http.StatusNotFound, Message: fmt.Sprintf("Role not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, role)
	}
}

