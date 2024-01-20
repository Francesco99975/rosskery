package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)


func AddProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.ProductDto
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		product, err := models.CreateProduct(payload.Name, payload.Description, payload.Price, payload.Image, payload.CategoryId, payload.Weighed)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error creating product: %v", err))
		}

		return c.JSON(http.StatusCreated, product)
	}
}


func Products() echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := models.GetProducts()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error while fetching products: %v", err))
		}

		return c.JSON(http.StatusOK, products)
	}
}


func Product() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Product
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		product, err := models.GetProduct(payload.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Product not found. Cause -> %v", err))
		}

		return c.JSON(http.StatusOK, product)
	}
}

func UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var payload models.ProductDto
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		product, err := models.GetProduct(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Product not found. Cause -> %v", err))
		}

		if err := product.Update(payload.Name, payload.Description, payload.Price, payload.Image, payload.Featured, payload.Published, payload.CategoryId, payload.Weighed); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error updating product: %v", err))
		}

		return c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		product, err := models.GetProduct(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Product not found. Cause -> %v", err))
		}

		defer func () {
			err = product.Delete()
			if err != nil {
				log.Errorf("Error deleting product: %v", err)
			}
		}()

		return c.JSON(http.StatusOK, product)
	}
}
