package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type OperationResult struct {
	value bool
}

type MessageUpdate struct {
	Message string `json:"message"`
}

func SetSetting(ctx context.Context, cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Setter

		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing request body for setting value: %v", err), Errors: []string{err.Error()}})
		}

		if err := storage.Valkey.Set(ctx, payload.Setting, payload.Value, 0).Err(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while setting value: %v", err), Errors: []string{err.Error()}})
		}

		online, err := storage.Valkey.Get(ctx, string(storage.Online)).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}
		operative, err := storage.Valkey.Get(ctx, string(storage.Operative)).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}
		message, err := storage.Valkey.Get(ctx, string(storage.Message)).Result()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}

		eventPayload, err := json.Marshal(storage.Settings{Message: message, Operative: operative, Online: online})
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while parsing payload: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventSettingsChanged, Payload: eventPayload})

		return c.JSON(http.StatusOK, payload)
	}
}

func GetSetting(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")

		val, err := storage.Valkey.Get(ctx, name).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}

		result := models.Setter{Setting: name, Value: val}

		log.Infof("GetSetting: %v", result)

		return c.JSON(http.StatusOK, result)
	}
}

func GetSettings(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		online, err := storage.Valkey.Get(ctx, string(storage.Online)).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}
		operative, err := storage.Valkey.Get(ctx, string(storage.Operative)).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}
		message, err := storage.Valkey.Get(ctx, string(storage.Message)).Result()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, storage.Settings{
			Online:    online,
			Operative: operative,
			Message:   message,
		})
	}
}

func SetMessage(ctx context.Context, cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload MessageUpdate

		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing request body for setting value: %v", err), Errors: []string{err.Error()}})
		}

		if err := storage.Valkey.Set(ctx, string(storage.Message), payload.Message, 0).Err(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while setting value: %v", err), Errors: []string{err.Error()}})
		}

		online, err := storage.Valkey.Get(ctx, string(storage.Online)).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}
		operative, err := storage.Valkey.Get(ctx, string(storage.Operative)).Bool()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}
		message, err := storage.Valkey.Get(ctx, string(storage.Message)).Result()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}

		eventPayload, err := json.Marshal(storage.Settings{Message: message, Operative: operative, Online: online})
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while parsing payload: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventSettingsChanged, Payload: eventPayload})

		return c.JSON(http.StatusOK, MessageUpdate{Message: message})
	}
}

func GetMessage(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		val, err := storage.Valkey.Get(ctx, string(storage.Message)).Result()

		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while getting setting: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, MessageUpdate{Message: val})
	}
}
