package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
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
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	e.Logger.SetLevel(log.INFO)
	e.GET("/healthcheck", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	e.Static("/assets", "./static")

	wsManager := models.NewManager(ctx)

	e.GET("/ws", wsManager.ServeWS)

	go wsManager.Run()

	e.GET("/", controllers.Index(ctx), middlewares.IsOnline(ctx))
	e.GET("/bag", controllers.GetCartItems(ctx))
	e.POST("/bag/:id", controllers.AddToCart(ctx))
	e.PUT("/bag/:id", controllers.RemoveOneFromCart(ctx))
	e.DELETE("/bag/:id", controllers.RemoveItemFromCart(ctx))
	e.DELETE("/bag", controllers.ClearCart(ctx))
	e.GET("/gallery", controllers.Gallery(ctx), middlewares.IsOnline(ctx))
	e.GET("/photos", controllers.Photos(), middlewares.IsOnline(ctx))
	e.GET("/shop", controllers.Shop(ctx), middlewares.IsOnline(ctx))

	e.POST("/login", api.Login(wsManager))
	e.POST("/check", api.CheckToken(wsManager))

	admin := e.Group("/admin")
	admin.Use(middlewares.IsAuthenticatedAdmin())
	admin.POST("/signup", api.Signup())
	admin.GET("/visits", api.GetVisits())
	admin.GET("/categories", api.Categories())
	admin.POST("/categories", api.CreateCategory())
	admin.DELETE("/categories/:id", api.DeleteCategory())
	admin.GET("/clientele", api.GetCustomerStats())
	admin.GET("/customers", api.Customers())
	admin.GET("/customers/:id", api.Customer())
	admin.DELETE("/customers/:id", api.DeleteCustomer())
	admin.GET("/finances", api.GetFinances())
	admin.GET("/orders", api.Orders())
	admin.GET("/orders/:id", api.Order())
	admin.POST("orders", api.IssueOrder())
	admin.DELETE("orders/:id", api.DeleteOrder())
	admin.GET("/products", api.Products())
	admin.GET("/products/:id", api.Product())
	admin.POST("/products", api.AddProduct())
	admin.PUT("/products/:id", api.UpdateProduct())
	admin.DELETE("/products/:id", api.DeleteProduct())
	admin.GET("/roles", api.Roles())
	admin.GET("/users", api.Users())
	admin.GET("/users/:id", api.User())
	admin.DELETE("/users/:id", api.DeleteUser())
	admin.GET("/setting", api.GetSetting(ctx))
	admin.PUT("/setting", api.SetSetting(ctx))
	admin.GET("/message", api.GetMessage(ctx))
	admin.PUT("/message", api.SetMessage(ctx))

	e.HTTPErrorHandler = serverErrorHandler

	return e
}

func serverErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	data := models.GetDefaultSite("Error", context.Background())

	buf := bytes.NewBuffer(nil)
	if code < 500 {
		_ = views.ClientError(data, err).Render(context.Background(), buf)

	} else {
		_ = views.ServerError(data, err).Render(context.Background(), buf)
	}

	_ = c.Blob(code, "text/html; charset=utf-8", buf.Bytes())

}
