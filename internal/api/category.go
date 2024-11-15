package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
)

func CreateCategory(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.CategoryDto
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for category: %v", err), Errors: []string{err.Error()}})
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error category not valid: %v", err), Errors: []string{err.Error()}})
		}

		categories, err := models.CreateCategory(payload.Category)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating category: %v", err), Errors: []string{err.Error()}})
		}

		newCategory, err := models.GetCategoryByName(payload.Category)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching category: %v", err), Errors: []string{err.Error()}})
		}

		rawNewCategory, err := json.Marshal(newCategory)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing category: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventNewCategory, Payload: rawNewCategory})

		return c.JSON(http.StatusCreated, categories)
	}
}

func Categories() echo.HandlerFunc {
	return func(c echo.Context) error {
		categories, err := models.GetCategories()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching categories: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, categories)
	}
}

func DeleteCategory(cm *models.ConnectionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		category, err := models.GetCategory(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching category while deleting: %v", err), Errors: []string{err.Error()}})
		}

		categories, err := category.Delete()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error deleting category: %v", err), Errors: []string{err.Error()}})
		}

		rawId, err := json.Marshal(struct{ Id string }{Id: id})
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing category: %v", err), Errors: []string{err.Error()}})
		}

		cm.BroadcastEvent(models.Event{Type: models.EventRemoveCategory, Payload: rawId})

		return c.JSON(http.StatusOK, categories)
	}
}
