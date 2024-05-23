package controllers

import (
	"context"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"
	"github.com/labstack/echo/v4"
)

func Shop(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Shop", ctx)

		products, err := models.GetPublishedProducts()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not fetch products")
		}

		html, err := helpers.GeneratePage(views.Shop(data, products))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
