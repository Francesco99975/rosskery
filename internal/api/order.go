package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)


func IssueOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		var payload models.OrderDto
		if err = c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for order: %v", err), Errors: []string{err.Error()}})
		}

		var customer *models.Customer

		if !models.CustomerExists(payload.Email) {
			customer, err = models.CreateCustomer(payload.Fullname, payload.Email, payload.Address, payload.Phone)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating customer: %v", err), Errors: []string{err.Error()}})
			}

		} else {
			customer, err = models.GetCustomerByEmail(payload.Email)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching customer: %v", err), Errors: []string{err.Error()}})
			}

			err = customer.Update(payload.Fullname, payload.Email, payload.Address, payload.Phone)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error updating customer: %v", err), Errors: []string{err.Error()}})
			}
		}

		order, err := models.CreateOrder(customer.Id, payload.Pickuptime, payload.PurchasedItems, payload.Method)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating order: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusCreated, order)
	}
}

func Orders() echo.HandlerFunc {
	return func(c echo.Context) error {
		orders, err := models.GetOrders()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching orders: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, orders)
	}
}

func Order() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching order: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, order)
	}
}

func DeleteOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{ Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching order while deleting: %v", err), Errors: []string{err.Error()}})
		}

		defer func () {
			err = order.Delete()
			if err != nil {
				log.Errorf("Error while deleting order: %v", err)
			}
		}()

		return c.JSON(http.StatusOK, order)
	}
}
