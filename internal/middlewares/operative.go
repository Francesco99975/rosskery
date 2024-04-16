package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/storage"
	"github.com/labstack/echo/v4"
)

func IsOperative(ctx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			operative, err := storage.Valkey.Get(ctx, string(storage.Operative)).Bool()

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("Settings error -> %v", err))
			}

			if !operative {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}

			return next(c)
		}
	}
}
