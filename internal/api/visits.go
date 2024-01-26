package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
)


func GetVisits() echo.HandlerFunc {
	return func (c echo.Context) error {
		visits, err := models.GetVisits()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{ Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching visits: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, visits)
	}
}
