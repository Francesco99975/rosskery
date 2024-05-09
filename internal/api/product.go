package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
)

func AddProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.ProductDto
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for product: %v", err), Errors: []string{err.Error()}})
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error product not valid: %v", err), Errors: []string{err.Error()}})
		}

		products, err := models.CreateProduct(payload.Name, payload.Description, payload.Price, payload.Image, payload.CategoryId, payload.Weighed)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating product: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusCreated, products)
	}
}

func Products() echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := models.GetProducts()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching products: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, products)
	}
}

func Product() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Product
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for product: %v", err), Errors: []string{err.Error()}})
		}

		product, err := models.GetProduct(payload.Id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error fetching product: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, product)
	}
}

func UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var payload models.ProductDto
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while updating product: %v", err), Errors: []string{err.Error()}})
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error product not valid: %v", err), Errors: []string{err.Error()}})
		}

		product, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Product not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		products, err := product.Update(payload.Name, payload.Description, payload.Price, payload.Image, payload.Featured, payload.Published, payload.CategoryId, payload.Weighed)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while updating product: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, products)
	}
}

func DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		product, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Product not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		products, err := product.Delete()
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error while deleting product: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, products)
	}
}
