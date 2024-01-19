package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
)

type ProductDto struct {
	name string
	description string
	price int
	image string
	featured bool
	published bool
	categoryId string
	weighed bool
}


func AddProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload ProductDto
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		product, err := models.CreateProduct(payload.name, payload.description, payload.price, payload.image, payload.categoryId, payload.weighed)
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
		var payload ProductDto
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Error parsing request body: %v", err))
		}

		product, err := models.GetProduct(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Product not found. Cause -> %v", err))
		}

		if err := product.Update(payload.name, payload.description, payload.price, payload.image, payload.featured, payload.published, payload.categoryId, payload.weighed); err != nil {
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

		if err := product.Delete(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Error deleting product: %v", err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}
