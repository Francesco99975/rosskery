package controllers

import (
	"context"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views"
	"github.com/labstack/echo/v4"
)

func PrivacyPolicy(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Privacy Policy", ctx)

		nonce := c.Get("nonce").(string)

		html, err := helpers.GeneratePage(views.PrivacyPolicy(data, nonce))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
