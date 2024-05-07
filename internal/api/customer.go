package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
)

func GetCustomerStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		timeframeStr := c.QueryParam("timeframe")

		timeframe := models.ParseTimeframe(timeframeStr)

		numberOfCustomers, err := models.GetAllCustomersAmount()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching customers amount: %v", err), Errors: []string{err.Error()}})
		}

		customerData, err := models.GetCustomersData(timeframe)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching customers data: %v", err), Errors: []string{err.Error()}})
		}

		topCustomers, err := models.GetTopSpenders()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching top spenders: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.CustomersStats{
			TotalCustomers: numberOfCustomers,
			CustomersData:  customerData,
			TopSpenders:    topCustomers,
		})

	}
}

func Customers() echo.HandlerFunc {
	return func(c echo.Context) error {
		customers, err := models.GetCustomers()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching customers: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, customers)
	}
}

func Customer() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		customer, err := models.GetCustomer(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching customer: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, customer)
	}
}

func DeleteCustomer() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		customer, err := models.GetCustomer(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching customer while deleting: %v", err), Errors: []string{err.Error()}})
		}

		customers, err := customer.Delete()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error deleting customer: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, customers)
	}
}
