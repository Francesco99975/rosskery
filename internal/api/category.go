package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)


type CategoryDto struct {
	category string
}

func CreateCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload CategoryDto
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing data for category: %v", err))
		}

		category, err := models.CreateCategory(payload.category)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error creating category: %v", err))
		}

		return c.JSON(http.StatusCreated, category)
	}
}

func Categories() echo.HandlerFunc {
	return func(c echo.Context) error {
		categories, err := models.GetCategories()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error fetching categories: %v", err))
		}

		return c.JSON(http.StatusOK, categories)
	}
}

func DeleteCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		category, err := models.GetCategory(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error fetching category while deleting: %v", err))
		}

		defer func () {
			err = category.Delete()
			if err != nil {
				log.Errorf("Error while deleting category: %v", err)
			}
		}()

		return c.JSON(http.StatusOK, category)
	}
}
