package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func AddProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		parsedPrice, err := strconv.Atoi(c.FormValue("price"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing price: %v", err), Errors: []string{err.Error()}})
		}

		parsedLv, err := strconv.Atoi(c.FormValue("lv"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing labor value: %v", err), Errors: []string{err.Error()}})
		}

		payload := models.ProductDto{
			Name:        c.FormValue("name"),
			Description: c.FormValue("description"),
			Price:       parsedPrice,
			CategoryId:  c.FormValue("category_id"),
			Weighed:     c.FormValue("weighed") == "true",
			Lv:          parsedLv,
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error product not valid: %v", err), Errors: []string{err.Error()}})
		}

		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		uploadedFiles := form.File["image"]
		if len(uploadedFiles) != 1 {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error uploading files: %v", err), Errors: nil})
		}

		file := uploadedFiles[0]

		products, err := models.CreateProduct(payload.Name, payload.Description, payload.Price, file, payload.CategoryId, payload.Weighed, payload.Lv)
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

		parsedPrice, err := strconv.Atoi(c.FormValue("price"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing price: %v", err), Errors: []string{err.Error()}})
		}

		parsedLv, err := strconv.Atoi(c.FormValue("lv"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing labor value: %v", err), Errors: []string{err.Error()}})
		}

		payload := models.ProductDto{
			Name:        c.FormValue("name"),
			Description: c.FormValue("description"),
			Price:       parsedPrice,
			CategoryId:  c.FormValue("category_id"),
			Weighed:     c.FormValue("weighed") == "true",
			Lv:          parsedLv,
			Published:   c.FormValue("published") == "true",
			Featured:    c.FormValue("featured") == "true",
		}

		if err := payload.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error product not valid: %v", err), Errors: []string{err.Error()}})
		}

		image := c.FormValue("file")
		var file *multipart.FileHeader
		log.Info("image: ", image)

		if !strings.Contains(image, "/assets/") {
			file, err = c.FormFile("image")
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error uploading files: %v", err), Errors: nil})
			}
		} else {
			file = nil
		}

		product, err := models.GetProduct(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Product not found. Cause -> %v", err), Errors: []string{err.Error()}})
		}

		products, err := product.Update(payload.Name, payload.Description, payload.Price, file, payload.Featured, payload.Published, payload.CategoryId, payload.Weighed, payload.Lv)
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
