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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching visits: %v", err))
		}

		return c.JSON(http.StatusOK, visits)
	}
}
