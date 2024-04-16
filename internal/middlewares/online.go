package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/helpers"

	"github.com/Francesco99975/rosskery/internal/storage"
	"github.com/Francesco99975/rosskery/views"
	"github.com/labstack/echo/v4"
)

func IsOnline(ctx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			online, err := storage.Valkey.Get(ctx, string(storage.Online)).Bool()

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Settings error -> %v", err))
			}

			if !online {
				html, err := helpers.GeneratePage(views.Offline())

				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page home")
				}

				return c.Blob(200, "text/html; charset=utf-8", html)
			}

			return next(c)
		}
	}
}
