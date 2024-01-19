package main

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/Francesco99975/rosskery/internal/api"
	"github.com/Francesco99975/rosskery/internal/controllers"
	"github.com/Francesco99975/rosskery/internal/middlewares"

	// "github.com/Francesco99975/rosskery/internal/middlewares"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func createRouter(ctx context.Context) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Logger.SetLevel(log.INFO)
	e.GET("/healthcheck", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	e.Static("/assets", "./static")

	wsManager := models.NewManager(ctx)

	e.GET("/ws", wsManager.ServeWS)

	go wsManager.Run()

	e.GET("/", controllers.Index())
	e.GET("/gallery", controllers.Gallery())
	e.GET("/photos", controllers.Photos())

	e.POST("/login", api.Login(wsManager))

	admin := e.Group("/admin")
	admin.Use(middlewares.IsAuthenticatedAdmin())
	admin.POST("/signup", api.Signup())



	e.HTTPErrorHandler = serverErrorHandler

	return e
}

func serverErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	data := models.GetDefaultSite("Error")

	buf := bytes.NewBuffer(nil)
	if code < 500 {
		_ = views.ClientError(data, err).Render(context.Background(), buf)

	} else {
		_ = views.ServerError(data, err).Render(context.Background(), buf)
	}

	_ = c.Blob(200, "text/html; charset=utf-8", buf.Bytes())

}
