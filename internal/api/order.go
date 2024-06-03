package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"
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

type OrderManager struct {
	cachedOrders map[string]models.OrderDto
	lock         sync.Mutex
}

var om = OrderManager{cachedOrders: make(map[string]models.OrderDto)}

func (o *OrderManager) Cache(id string, payload models.OrderDto) {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.cachedOrders[id] = payload
}

func (o *OrderManager) Confirm(ctx context.Context, c echo.Context, id string) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	payload := o.cachedOrders[id]
	if err := processOrder(ctx, c, payload); err != nil {
		return err
	}

	delete(o.cachedOrders, id)
	return nil
}

func processOrder(ctx context.Context, c echo.Context, payload models.OrderDto) error {
	var err error

	var customer *models.DbCustomer
	exists, err := models.CustomerExists(payload.Email)
	if err != nil {
		return fmt.Errorf("Error checking if customer exists: %v", err)
	}

	if !exists {
		customer, err = models.CreateCustomer(payload.Fullname, payload.Email, payload.Address, payload.Phone)
		if err != nil {
			return fmt.Errorf("Error creating customer: %v", err)
		}

	} else {
		customer, err = models.GetCustomerByEmail(payload.Email)
		if err != nil {
			return fmt.Errorf("Error fetching customer: %v", err)
		}

		err := customer.Update(payload.Fullname, payload.Email, payload.Address, payload.Phone)
		if err != nil {
			return fmt.Errorf("Error updating customer: %v", err)
		}
	}

	sess, err := session.Get("session", c)

	if err != nil {
		return fmt.Errorf("Error fetching session: %v", err)
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		// Secure:   true,
		// Domain:   "",
		// SameSite: http.SameSiteDefaultMode,
	}
	sessionID, ok := sess.Values["sessionID"].(string)
	if !ok || sessionID == "" {

		return errors.New("Session ID is invalid")

	}
	cart, err := models.GetCart(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("Error fetching cart: %v", err)
	}

	purchases, err := cart.Purchases()
	if err != nil {
		return fmt.Errorf("Error fetching purchases: %v", err)
	}

	_, err = models.CreateOrder(customer.Id, payload.Pickuptime, purchases, payload.Method)
	if err != nil {
		return fmt.Errorf("Error creating order: %v", err)
	}

	if err = cart.Clear(ctx); err != nil {
		return fmt.Errorf("Error clearing cart: %v", err)
	}

	return nil
}

func PaymentWebhook(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			// Secure:   true,
			// Domain:   "",
			// SameSite: http.SameSiteDefaultMode,
		}
		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {

			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create session")

		}

		const MaxBodyBytes = int64(65536)
		body := http.MaxBytesReader(c.Response().Writer, c.Request().Body, MaxBodyBytes)
		payload, err := io.ReadAll(body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error reading body")
		}

		event := stripe.Event{}
		if err := json.Unmarshal(payload, &event); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error parsing request body")
		}

		switch event.Type {
		case "payment_intent.succeeded":
			if err := om.Confirm(ctx, c, sessionID); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Error confirming order")
			}
			data := models.GetDefaultSite("Order Confirmed", ctx)

			html, err := helpers.GeneratePage(views.Confirmation(data))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
			}

			return c.Blob(200, "text/html; charset=utf-8", html)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "Unhandled event type")
		}

	}
}

func IssueOrder(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		date, err := time.Parse("2006-01-02T15:04", c.FormValue("pickuptime"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error validating order at payment check: %v", err), Errors: []string{err.Error()}})
		}

		payload := models.OrderDto{
			Email:      c.FormValue("email"),
			Fullname:   c.FormValue("fullname"),
			Phone:      c.FormValue("phone"),
			Address:    c.FormValue("address"),
			Pickuptime: date,
			Method:     models.ParsePaymentMethod(c.FormValue("method")),
		}

		if err = payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error validating order at payment check: %v", err), Errors: []string{err.Error()}})
		}

		if payload.Method == models.CASH {
			if err := processOrder(ctx, c, payload); err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error processing order: %v", err), Errors: []string{err.Error()}})
			}

			data := models.GetDefaultSite("Order Confirmed", ctx)

			html, err := helpers.GeneratePage(views.Confirmation(data))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
			}

			return c.Blob(200, "text/html; charset=utf-8", html)

		}

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			// Secure:   true,
			// Domain:   "",
			// SameSite: http.SameSiteDefaultMode,
		}
		sessionID, ok := sess.Values["sessionID"].(string)
		if !ok || sessionID == "" {

			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create session")

		}

		om.Cache(sessionID, payload)

		data := models.GetDefaultSite("Pay Online", ctx)

		html, err := helpers.GeneratePage(views.Pay(data, os.Getenv("STRIPE_PUBLISHABLE_KEY")))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
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

		orders, err := order.Delete()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error deleting order: %v", err), Errors: []string{err.Error()}})
		}
		return c.JSON(http.StatusOK, orders)
	}
}
