package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Francesco99975/rosskery/internal/api"
	"github.com/Francesco99975/rosskery/internal/controllers"
	"github.com/Francesco99975/rosskery/internal/middlewares"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"

	"github.com/gorilla/sessions"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func createRouter(ctx context.Context) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	e.Use(middlewares.RateLimiter)

	e.Use(middlewares.BrotliMiddleware)

	e.Use(middleware.Gzip())

	e.GET("/healthcheck", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	e.Static("/assets", "./static")

	wsManager := models.NewManager(ctx)

	e.GET("/ws", wsManager.ServeWS)

	web := e.Group("")

	web.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:_csrf,header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookiePath:     "/",
		CookieHTTPOnly: true,
	}))

	if os.Getenv("GO_ENV") == "development" {
		e.Logger.SetLevel(log.DEBUG)
		web.Use(middlewares.SecurityHeadersDev())
	}

	if os.Getenv("GO_ENV") == "production" {
		e.Logger.SetLevel(log.WARN)
		web.Use(middlewares.SecurityHeaders())
	}

	go wsManager.Run()

	web.GET("/", controllers.Index(ctx), middlewares.IsOnline(ctx))
	web.GET("/policy", controllers.PrivacyPolicy(ctx), middlewares.IsOnline(ctx))
	web.GET("/terms", controllers.Terms(ctx), middlewares.IsOnline(ctx))
	web.GET("/gallery", controllers.Gallery(ctx), middlewares.IsOnline(ctx))
	web.GET("/photos", controllers.Photos(), middlewares.IsOnline(ctx))
	web.GET("/shop", controllers.Shop(ctx), middlewares.IsOnline(ctx))
	web.GET("/checkout", controllers.Checkout(ctx), middlewares.IsOnline(ctx), middlewares.IsOperative(ctx))

	web.GET("/bag", controllers.GetCartItems(ctx), middlewares.IsOnline(ctx))
	web.POST("/bag/:id", controllers.AddToCart(ctx), middlewares.IsOnline(ctx))
	web.PUT("/bag/:id", controllers.RemoveOneFromCart(ctx), middlewares.IsOnline(ctx))
	web.DELETE("/bag/:id", controllers.RemoveItemFromCart(ctx), middlewares.IsOnline(ctx))
	web.DELETE("/bag", controllers.ClearCart(ctx), middlewares.IsOnline(ctx))
	web.POST("/intent", api.CreatePaymentIntent(ctx), middlewares.IsOnline(ctx))
	web.POST("/orders", api.IssueOrder(ctx, wsManager), middlewares.IsOnline(ctx))
	web.GET("/orders/success", controllers.Success(ctx), middlewares.IsOnline(ctx))

	web.GET("/address", controllers.AddressAutocomplete())

	web.POST("/webhook", api.PaymentWebhook(ctx, wsManager))

	admin := e.Group("/admin")
	admin.POST("/login", api.Login(wsManager))
	admin.POST("/check", api.CheckToken(wsManager))
	// admin.GET("/notif", api.GetNotifConfig())

	admin.Use(middlewares.IsAuthenticatedAdmin())

	admin.POST("/signup", api.Signup())
	admin.GET("/visits", api.GetVisits())
	admin.GET("/visits/stats", api.GetVisitStats())
	admin.GET("/visits/graph", api.GetVisitGraph())
	admin.GET("/visits/standings", api.GetVisitsStandings())
	admin.GET("/categories", api.Categories())
	admin.POST("/categories", api.CreateCategory(wsManager))
	admin.DELETE("/categories/:id", api.DeleteCategory(wsManager))
	admin.GET("/clientele", api.GetCustomerStats())
	admin.GET("/customers", api.Customers())
	admin.GET("/customers/:id", api.Customer())
	admin.DELETE("/customers/:id", api.DeleteCustomer(wsManager))
	admin.GET("/finances", api.GetFinances())
	admin.GET("/finances/stats", api.GetFinancesStats())
	admin.GET("/finances/orders", api.GetOrdersData())
	admin.GET("/finances/monetary", api.GetMonetaryData())
	admin.GET("/finances/payments", api.GetPaymentData())
	admin.GET("/finances/status", api.GetOrdersStatusPie())
	admin.GET("/finances/methods", api.GetOrdersPaymentPie())
	admin.GET("/finances/standings", api.GetOrdersStandings())
	admin.GET("/orders", api.Orders())
	admin.GET("/orders/:id", api.Order())
	admin.GET("/fulfill/:id", api.FulfillOrder())
	// admin.POST("orders", api.IssueOrder(ctx))
	admin.DELETE("orders/:id", api.DeleteOrder())
	admin.GET("/products", api.Products())
	admin.GET("/products/:id", api.Product())
	admin.POST("/products", api.AddProduct(wsManager))
	admin.PUT("/products/:id", api.UpdateProduct(wsManager))
	admin.DELETE("/products/:id", api.DeleteProduct(wsManager))
	admin.GET("/roles", api.Roles())
	admin.GET("/users", api.Users())
	admin.GET("/users/:id", api.User())
	admin.DELETE("/users/:id", api.DeleteUser())
	admin.GET("/setting/:name", api.GetSetting(ctx))
	admin.PUT("/setting", api.SetSetting(ctx, wsManager))
	admin.GET("/message", api.GetMessage(ctx))
	admin.PUT("/message", api.SetMessage(ctx, wsManager))

	e.HTTPErrorHandler = serverErrorHandler

	return e
}

func serverErrorHandler(err error, c echo.Context) {
	// Default to internal server error (500)
	code := http.StatusInternalServerError
	var message interface{} = "An unexpected error occurred"

	// Check if it's an echo.HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	}

	// Check the Accept header to decide the response format
	if strings.Contains(c.Request().Header.Get("Accept"), "application/json") {
		// Respond with JSON if the client prefers JSON
		_ = c.JSON(code, map[string]interface{}{
			"error":   true,
			"message": message,
			"status":  code,
		})
	} else {
		// Prepare data for rendering the error page (HTML)
		data := models.GetDefaultSite("Error", context.Background())

		// Buffer to hold the HTML content (in case of HTML response)
		buf := bytes.NewBuffer(nil)

		// Render based on the status code
		if code >= 500 {
			_ = views.ServerError(data, err).Render(context.Background(), buf)
		} else {
			_ = views.ClientError(data, err).Render(context.Background(), buf)
		}
		// Respond with HTML (default) if the client prefers HTML
		_ = c.Blob(code, "text/html; charset=utf-8", buf.Bytes())
	}
}
