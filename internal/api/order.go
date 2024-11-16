package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/internal/tools"
	"github.com/Francesco99975/rosskery/views"
	"github.com/Francesco99975/rosskery/views/components"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

func GetOrdersData() echo.HandlerFunc {
	return func(c echo.Context) error {
		timeframeStr := c.QueryParam("timeframe")
		methodStr := c.QueryParam("method")
		status := c.QueryParam("status") == "true"

		timeframe := models.ParseTimeframe(timeframeStr)
		method := models.ParsePaymentMethod(methodStr)

		ordersData, err := models.GetOrdersData(timeframe, method, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching orders data: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.Graph{Data: ordersData})
	}
}

func GetMonetaryData() echo.HandlerFunc {
	return func(c echo.Context) error {
		timeframeStr := c.QueryParam("timeframe")
		methodStr := c.QueryParam("method")
		status := c.QueryParam("status") == "true"

		timeframe := models.ParseTimeframe(timeframeStr)
		method := models.ParsePaymentMethod(methodStr)

		monetaryData, err := models.GetMonetaryData(timeframe, method, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching monetary data: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.Graph{Data: monetaryData})
	}
}

func GetPaymentData() echo.HandlerFunc {
	return func(c echo.Context) error {
		timeframeStr := c.QueryParam("timeframe")
		status := c.QueryParam("status") == "true"

		timeframe := models.ParseTimeframe(timeframeStr)

		preferredMethodData, err := models.GetPreferredMethodData(timeframe, status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching preferred method data: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.Graphs{Datapoints: preferredMethodData})
	}
}

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

func GetFinancesStats() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		return c.JSON(http.StatusOK, models.FinancesStats{OrdersAmount: numberOfOrders, OutstandingCash: outstanding, PendingMoney: pending, Gains: gains, Total: total})
	}
}

func GetOrdersStatusPie() echo.HandlerFunc {
	return func(c echo.Context) error {
		filledPie, err := models.GetFilledPie()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching filled pie: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.PieGraph{Pie: filledPie})
	}
}

func GetOrdersPaymentPie() echo.HandlerFunc {
	return func(c echo.Context) error {
		paymentMethodPie, err := models.GetMethodsPie()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching payment method pie: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.PieGraph{Pie: paymentMethodPie})
	}
}

func GetOrdersStandings() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		return c.JSON(http.StatusOK, models.OrdersStandingsResponse{RankedOrders: topOrders, ToppedSellers: topSellers, FloppedSellers: flopSellers, ToppedGainers: topGainers, FloppedGainers: flopGainers})
	}
}

type OrderManager struct {
	cachedOrders        map[string]models.OrderDto
	realtedCreationDate map[string]time.Time
	lock                sync.Mutex
}

func NewOrderManager() *OrderManager {
	om := &OrderManager{cachedOrders: make(map[string]models.OrderDto, 0), realtedCreationDate: make(map[string]time.Time, 0)}
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			if err := om.AutoClean(); err != nil {
				log.Printf("Error cleaning orders: %v", err)
			}
		}
	}()
	return om
}

var om = NewOrderManager()

func (o *OrderManager) Cache(id string, payload models.OrderDto) {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.cachedOrders[id] = payload
	o.realtedCreationDate[id] = time.Now()
}

func (o *OrderManager) Confirm(ctx context.Context, id string, cm *models.ConnectionManager) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	payload := o.cachedOrders[id]
	if err := processOrder(ctx, payload, id, cm); err != nil {
		return err
	}

	delete(o.cachedOrders, id)
	delete(o.realtedCreationDate, id)
	return nil
}

func (o *OrderManager) AutoClean() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	for id, creationDate := range o.realtedCreationDate {
		if time.Since(creationDate) > 10*time.Minute {
			delete(o.cachedOrders, id)
			delete(o.realtedCreationDate, id)
		}
	}

	return nil
}

func processOrder(ctx context.Context, payload models.OrderDto, sessionID string, cm *models.ConnectionManager) error {
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

	cart, err := models.GetCart(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("Error fetching cart: %v", err)
	}

	purchases, err := cart.Purchases()
	if err != nil {
		return fmt.Errorf("Error fetching purchases: %v", err)
	}

	order, err := models.CreateOrder(customer.Id, payload.Pickuptime, purchases, payload.Method)
	if err != nil {
		return fmt.Errorf("Error creating order: %v", err)
	}

	if err = cart.Clear(ctx); err != nil {
		return fmt.Errorf("Error clearing cart: %v", err)
	}

	total := helpers.FormatPrice(float64(helpers.FoldSlice[models.Purchase, func(models.Purchase, int) int, int](order.Purchases, func(prev models.Purchase, cur int) int {
		return prev.Product.Price*prev.Quantity + cur
	}, 0)) / 100.0)

	invoice, err := tools.GenerateInvoice(order)
	if err != nil {
		return fmt.Errorf("Error generating invoice: %v", err)
	}

	payStatus := "Pay at Pickup"
	if models.ParsePaymentMethod(order.Method) != models.CASH {
		payStatus = "No payment is due"
	}

	purchaseDetails := helpers.MapSlice[models.Purchase, tools.ReceiptDetail](order.Purchases, func(p models.Purchase) tools.ReceiptDetail {
		return tools.ReceiptDetail{Description: fmt.Sprintf("%s - (x%d)", p.Product.Name, p.Quantity), Amount: helpers.FormatPrice(float64(p.Product.Price*p.Quantity) / 100.0)}
	})

	err = tools.SendReceipt(order.Customer.Email, tools.Receipt{ProductURL: "rosskery.com", ProductName: "Rosskery", Customer: order.Customer.Fullname, PaymentStatus: payStatus, CreditCardStatementName: "Rosskery", OrderID: order.Id, Date: order.Created.Format("2006-01-02 03:04 PM"), PickupDate: order.Pickuptime.Format("2006-01-02 03:04 PM"), ReceiptDetails: purchaseDetails, Total: fmt.Sprint(total), SupportURL: "", CompanyName: "Rosskey", CompanyAddress: "robarra@rosskery.com"}, invoice)
	if err != nil {
		return fmt.Errorf("Error sending receipt: %v", err)
	}

	cm.BroadcastEvent(models.Event{Type: models.EventOrdersChanged, Payload: nil})
	cm.BroadcastEvent(models.Event{Type: models.EventCustomersChanged, Payload: nil})

	return nil
}

func PaymentWebhook(ctx context.Context, cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		const MaxBodyBytes = int64(65536)
		body := http.MaxBytesReader(c.Response().Writer, c.Request().Body, MaxBodyBytes)
		payload, err := io.ReadAll(body)
		if err != nil {
			log.Errorf("Error reading body: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Error reading body")
		}

		event := stripe.Event{}
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Errorf("Error parsing request body: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Error parsing request body")
		}

		signatureHeader := c.Request().Header.Get("Stripe-Signature")
		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

		event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
		if err != nil {
			log.Errorf("Error constructing event: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Error constructing event")
		}

		switch event.Type {
		case "payment_intent.succeeded":
			var paymentIntent stripe.PaymentIntent
			err := json.Unmarshal(event.Data.Raw, &paymentIntent)
			if err != nil {
				log.Errorf("Error parsing payment intent: %v", err)
				return echo.NewHTTPError(http.StatusBadRequest, "Error parsing payment intent")
			}

			sessionID := paymentIntent.Metadata["sessionID"]
			if err := om.Confirm(ctx, sessionID, cm); err != nil {
				log.Errorf("Error confirming order: %v", err)
				return echo.NewHTTPError(http.StatusBadRequest, "Error confirming order")
			}

			data := models.GetDefaultSite("Order Confirmed", ctx)
			nonce := c.Get("nonce").(string)

			html, err := helpers.GeneratePage(views.Confirmation(data, nonce))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
			}

			return c.Blob(200, "text/html; charset=utf-8", html)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "Unhandled event type")
		}

	}
}

func IssueOrder(ctx context.Context, cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		date, err := time.Parse("2006-01-02 15:04", c.FormValue("pickuptime"))
		if err != nil {
			log.Errorf("Error parsing pickuptime: %v", err)
			html, err := helpers.GeneratePage(components.Errors("Invalid date for pickuptime"))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
			}

			return c.Blob(http.StatusBadRequest, "text/html; charset=utf-8", html)
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
			log.Errorf("Error validating payload: %v", err)
			html, err := helpers.GeneratePage(components.Errors(fmt.Sprintf("Error: %v", err)))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
			}

			return c.Blob(http.StatusBadRequest, "text/html; charset=utf-8", html)
		}

		log.Infof("Payload: %v", payload)

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

		if payload.Method == models.CASH {
			if err := processOrder(ctx, payload, sessionID, cm); err != nil {
				log.Errorf("Error processing order: %v", err)
				html, err := helpers.GeneratePage(components.Errors("Error processing order"))
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
				}

				return c.Blob(http.StatusBadRequest, "text/html; charset=utf-8", html)
			}

			data := models.GetDefaultSite("Order Confirmed", ctx)
			nonce := c.Get("nonce").(string)

			html, err := helpers.GeneratePage(views.Confirmation(data, nonce))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
			}

			return c.Blob(200, "text/html; charset=utf-8", html)

		}

		om.Cache(sessionID, payload)

		data := models.GetDefaultSite("Pay Online", ctx)

		csrfToken := c.Get("csrf").(string)
		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(views.Pay(data, os.Getenv("STRIPE_PUBLISHABLE_KEY"), csrfToken, nonce))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}

func FulfillOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		order, err := models.GetOrder(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching order while fulfilling: %v", err), Errors: []string{err.Error()}})
		}

		updatedOrder, err := order.Fulfill()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fulfilling order: %v", err), Errors: []string{err.Error()}})
		}
		return c.JSON(http.StatusOK, updatedOrder)

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
