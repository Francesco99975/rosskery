package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/internal/storage"
	"github.com/labstack/echo/v4"
)

type OperationResult struct {
	value bool
}

func SetSetting(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Setter

		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing request body for setting value: %v", err), Errors: []string{err.Error()}})
		}

		if err := storage.Valkey.Set(ctx, payload.Setting, payload.Value, 0).Err(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while setting value: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, OperationResult{value: payload.Value})
	}
}

func GetSetting(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Setter

		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing request body getting setting: %v", err), Errors: []string{err.Error()}})
		}

		val, err := storage.Valkey.Get(ctx, payload.Setting).Bool()

		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, OperationResult{value: val})
	}
}
