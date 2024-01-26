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
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		var customer *models.Customer

		if !models.CustomerExists(payload.Email) {
			customer, err = models.CreateCustomer(payload.Fullname, payload.Email, payload.Address, payload.Phone)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error creating customer: %v", err))
			}

		} else {
			customer, err = models.GetCustomerByEmail(payload.Email)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error creating customer: %v", err))
			}

			err = customer.Update(payload.Fullname, payload.Email, payload.Address, payload.Phone)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error updating customer: %v", err))
			}
		}

		order, err := models.CreateOrder(customer.Id, payload.Pickuptime, payload.PurchasedItems, payload.Method)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error creating order: %v", err))
		}

		return c.JSON(http.StatusCreated, order)
	}
}

func Orders() echo.HandlerFunc {
	return func(c echo.Context) error {
		orders, err := models.GetOrders()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching orders: %v", err))
		}

		return c.JSON(http.StatusOK, orders)
	}
}

func Order() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching order: %v", err))
		}

		return c.JSON(http.StatusOK, order)
	}
}

func DeleteOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching order: %v", err))
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
