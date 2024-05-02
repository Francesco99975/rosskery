package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GetFinances() echo.HandlerFunc {
	return func(c echo.Context) error {

		timeframeStr := c.QueryParam("timeframe")
		methodStr := c.QueryParam("method")
		status := c.QueryParam("status") == "true"

		timeframe := models.ParseTimeframe(timeframeStr)
		method := models.ParsePaymentMethod(methodStr)

		numberOfOrders, err := models.GetOrdersAmount()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching orders amount: %v", err), Errors: []string{err.Error()}})
		}

		outstanding, err := models.GetOutstandingCash()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching outstanding cash: %v", err), Errors: []string{err.Error()}})
		}

		pending, err := models.GetPendingMoney()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching pending money: %v", err), Errors: []string{err.Error()}})
		}

		gains, err := models.GetGains()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching gains: %v", err), Errors: []string{err.Error()}})
		}

		total, err := models.GetTotalFromOrders()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching total from orders: %v", err), Errors: []string{err.Error()}})
		}

		ordersData, err := models.GetOrdersData(timeframe, method, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching orders data: %v", err), Errors: []string{err.Error()}})
		}

		monetaryData, err := models.GetMonetaryData(timeframe, method, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching monetary data: %v", err), Errors: []string{err.Error()}})
		}

		preferredMethodData, err := models.GetPreferredMethodData(timeframe, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching preferred method data: %v", err), Errors: []string{err.Error()}})
		}

		filledPie, err := models.GetFilledPie()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching filled pie: %v", err), Errors: []string{err.Error()}})
		}

		paymentMethodPie, err := models.GetMethodsPie()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching payment method pie: %v", err), Errors: []string{err.Error()}})
		}

		topOrders, err := models.GetTopOrders()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching top orders: %v", err), Errors: []string{err.Error()}})
		}

		topSellers, err := models.GetTopSellers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching top sellers: %v", err), Errors: []string{err.Error()}})
		}

		flopSellers, err := models.GetFlopSellers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching flop sellers: %v", err), Errors: []string{err.Error()}})
		}

		topGainers, err := models.GetTopGainers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching top gainers: %v", err), Errors: []string{err.Error()}})
		}

		flopGainers, err := models.GetFlopGainers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching flop gainers: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.FinancesResponse{OrdersAmount: numberOfOrders, OutstandingCash: outstanding, PendingMoney: pending, Gains: gains, Total: total, OrdersData: ordersData, MonetaryData: monetaryData, PreferredMethodData: preferredMethodData, FilledPie: filledPie, MethodPie: paymentMethodPie, RankedOrders: topOrders, ToppedSellers: topSellers, FloppedSellers: flopSellers, ToppedGainers: topGainers, FloppedGainers: flopGainers})
	}
}

func IssueOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		var payload models.OrderDto
		if err = c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for order: %v", err), Errors: []string{err.Error()}})
		}

		var customer *models.Customer

		if !models.CustomerExists(payload.Email) {
			customer, err = models.CreateCustomer(payload.Fullname, payload.Email, payload.Address, payload.Phone)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating customer: %v", err), Errors: []string{err.Error()}})
			}

		} else {
			customer, err = models.GetCustomerByEmail(payload.Email)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching customer: %v", err), Errors: []string{err.Error()}})
			}

			err = customer.Update(payload.Fullname, payload.Email, payload.Address, payload.Phone)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error updating customer: %v", err), Errors: []string{err.Error()}})
			}
		}

		order, err := models.CreateOrder(customer.Id, payload.Pickuptime, payload.PurchasedItems, payload.Method)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating order: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusCreated, order)
	}
}

func Orders() echo.HandlerFunc {
	return func(c echo.Context) error {
		orders, err := models.GetOrders()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching orders: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, orders)
	}
}

func Order() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching order: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, order)
	}
}

func DeleteOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching order while deleting: %v", err), Errors: []string{err.Error()}})
		}

		defer func() {
			err = order.Delete()
			if err != nil {
				log.Errorf("Error while deleting order: %v", err)
			}
		}()

		return c.JSON(http.StatusOK, order)
	}
}
