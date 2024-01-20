package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Customers() echo.HandlerFunc {
	return func(c echo.Context) error {
		customers, err := models.GetCustomers()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error fetching customers: %v", err))
		}

		return c.JSON(http.StatusOK, customers)
	}
}

func Customer() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		customer, err := models.GetCustomer(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error fetching customer: %v", err))
		}

		return c.JSON(http.StatusOK, customer)
	}
}


func DeleteCustomer() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		customer, err := models.GetCustomer(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error fetching customer while deleting: %v", err))
		}

		defer func () {
			err = customer.Delete()
			if err != nil {
				log.Errorf("Error while deleting customer: %v", err)
			}
		}()

		return c.JSON(http.StatusOK, customer)
	}
}

