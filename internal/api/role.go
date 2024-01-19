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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching roles: %v", err))
		}

		return c.JSON(http.StatusOK, roles)
	}
}

func Role() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Role
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		role, err := models.GetRoleById(payload.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Role not found. Cause -> %v", err))
		}

		return c.JSON(http.StatusOK, role)
	}
}

